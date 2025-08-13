package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mrizkifadil26/medix/logger"
)

type Walker struct {
	Context  context.Context
	Opts     WalkOptions
	Stats    *WalkStats
	Progress *ProgressTracker
	Logger   logger.Logger

	OnVisitFile func(path string, size int64) error
	OnVisitDir  func(path string, entries []fs.DirEntry) error
	OnSkip      func(path string, reason string) error
	OnError     func(path string, err error) error

	mu sync.Mutex
}

type WalkOptions struct {
	StopOnError bool
	SkipOnError bool

	OnlyDirs        bool
	OnlyFiles       bool
	MaxDepth        int
	MinIncludeDepth int
	SkipEmptyDirs   bool
	SkipRoot        bool
	OnlyLeafDirs    bool
	IncludeHidden   bool

	IncludePatterns []string
	ExcludePatterns []string
	IncludeExts     []string
	ExcludeExts     []string

	IncludeErrors bool
	IncludeStats  bool // collect walk statistics

	EnableProgress bool

	Debug DebugOptions
}

// DebugOptions holds debug logging settings.
type DebugOptions struct {
	Enable bool
	Level  string // e.g. "ERROR", "DEBUG", "TRACE"
}

type WalkStats struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration // EndTime - StartTime

	EntriesVisited int // total entries seen (files + dirs + others)
	FilesVisited   int
	DirsVisited    int
	Matches        int // entries that passed all filters
	Skipped        int // entries skipped by filters or prune
	ErrorsCount    int
	// SymlinksVisited int

	// --- Size & Resource Tracking ---
	TotalSize   int64 // bytes processed
	AvgFileSize int64
	MinFileSize int64
	MaxFileSize int64
	DataRate    float64 // bytes/sec

	Errors []error

	EntriesPerSec float64 // processing rate
	FilesPerSec   float64

	// --- Custom User Metrics ---
	Custom map[string]interface{}
}

// Default constructor without logger (logger will be nil)
func NewWalker(ctx context.Context, opts WalkOptions) *Walker {
	return NewWalkerWithLogger(ctx, opts, nil)
}

func NewWalkerWithLogger(ctx context.Context, opts WalkOptions, log logger.Logger) *Walker {
	var loggerWithCtx logger.Logger

	if log == nil {
		// Create default SimpleLogger
		defaultLogger := logger.NewLogrusLogger()
		loggerWithCtx = defaultLogger.WithContext("scanner.Walker")
	} else {
		// Use provided logger, ensure output set
		loggerWithCtx = log.WithContext("scanner.Walker")
	}

	w := &Walker{
		Context: ctx,
		Opts:    opts,
		Stats:   &WalkStats{Custom: make(map[string]interface{})},
		Logger:  loggerWithCtx,
	}

	// Configure logger based on debug opts
	if w.Logger != nil && opts.Debug.Enable {
		level := parseLevel(opts.Debug.Level)

		w.Logger.SetLevel(level)
		w.Logger.SetEnabled(true)
	} else if w.Logger != nil {
		w.Logger.SetEnabled(false)
	}

	return w
}

func (w *Walker) Walk(root string) error {
	start := time.Now()
	w.debug(root, "Starting walk", map[string]interface{}{
		"IncludeStats":    w.Opts.IncludeStats,
		"EnableProgress":  w.Opts.EnableProgress,
		"MaxDepth":        w.Opts.MaxDepth,
		"MinIncludeDepth": w.Opts.MinIncludeDepth,
	})

	// Initialize stats if enabled
	if w.Opts.IncludeStats && w.Stats != nil {
		w.debug(root, "Initializing stats", nil)

		w.Stats.StartTime = start
		w.Stats.MinFileSize = -1 // sentinel for tracking min
		if w.Stats.Custom == nil {
			w.Stats.Custom = make(map[string]interface{})
		}

		if w.Stats.Errors == nil {
			w.Stats.Errors = []error{}
		}
	}

	// Initialize progress if enabled
	if w.Opts.EnableProgress {
		w.debug(root, "Counting total entries for progress", nil)

		// Count total entries first (no callbacks)
		progressStats, _ := w.Count(root)
		total := progressStats.EntriesVisited
		w.debug(root, "Total entries found", map[string]interface{}{"total": total})

		if total > 0 {
			w.Progress = NewProgressTracker(total, true, "Scanning")
			defer w.Progress.Finish()
		}
	}

	defer func() {
		if w.Opts.IncludeStats && w.Stats != nil {
			w.mu.Lock()
			w.Stats.EndTime = time.Now()
			w.Stats.Duration = w.Stats.EndTime.Sub(w.Stats.StartTime)

			// Calculate derived stats
			if w.Stats.FilesVisited > 0 {
				w.Stats.AvgFileSize = w.Stats.TotalSize / int64(w.Stats.FilesVisited)
			}

			if w.Stats.Duration > 0 {
				durSec := w.Stats.Duration.Seconds()
				w.Stats.EntriesPerSec = float64(w.Stats.EntriesVisited) / durSec
				w.Stats.FilesPerSec = float64(w.Stats.FilesVisited) / durSec
				w.Stats.DataRate = float64(w.Stats.TotalSize) / durSec
			}

			w.mu.Unlock()

			stats := w.GetStats()
			w.debug(root, "Stats collected", FormatStats(stats))
		}
	}()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			w.error(path, "WalkDir callback error", map[string]interface{}{"error": err})
			return w.handleError(path, err)
		}

		// Respect context cancellation
		select {
		case <-w.Context.Done():
			w.debug(path, "Context canceled", nil)
			w.handleSkip(path, "context canceled")
			return w.handleError(path, w.Context.Err())
		default:
		}

		depth := getDepth(root, path)
		w.trace(path, "Visiting entry", map[string]interface{}{
			"depth": depth,
			"isDir": d.IsDir(),
			"name":  d.Name(),
		})

		// Skip entries above MinIncludeDepth from being added to results/stats/callbacks
		includeThis := depth >= w.Opts.MinIncludeDepth

		// --- Stats: entries visited ---
		if includeThis && w.Opts.IncludeStats && w.Stats != nil {
			w.mu.Lock()
			w.Stats.EntriesVisited++
			w.mu.Unlock()
		}

		// Hidden file/dir skipping
		if !w.Opts.IncludeHidden && isHidden(d.Name()) {
			w.trace(path, "Skipping hidden entry", nil)
			return w.handleSkip(path, "hidden entry")
		}

		if w.Opts.SkipRoot && depth == 0 {
			w.debug(path, "Skipping root directory itself", nil)
			return nil
		}

		if w.Opts.MaxDepth >= 0 && depth > w.Opts.MaxDepth {
			w.trace(path, "Skipping due to max depth limit", map[string]interface{}{
				"maxDepth": w.Opts.MaxDepth,
				"current":  depth,
			})

			return fs.SkipDir
		}

		// Skip based on OnlyDirs / OnlyFiles flags:
		if w.Opts.OnlyDirs && !d.IsDir() {
			w.trace(path, "Skipping file due to OnlyDirs flag", nil)
			return nil // skip files
		}
		if w.Opts.OnlyFiles && d.IsDir() {
			w.trace(path, "Skipping dir due to OnlyFiles flag", nil)
			return nil // skip dirs
		}

		// --- Progress tracking ---
		if includeThis && w.Opts.EnableProgress && w.Progress != nil {
			w.Progress.Increment(1)
		}

		// Directory handling
		if d.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				w.error(path, "Error reading directory", map[string]interface{}{"err": err})
				return w.handleError(path, err)
			}

			w.trace(path, "Read directory entries", map[string]interface{}{"entriesCount": len(entries)})

			// Skip empty dirs
			if w.Opts.SkipEmptyDirs && len(entries) == 0 {
				w.trace(path, "Skipping empty directory", nil)
				return w.handleSkip(path, "empty dir")
			}

			// Only leaf dirs (skip if not leaf)
			if w.Opts.OnlyLeafDirs && d.IsDir() && !isLeafDir(path) {
				return nil
			}

			// Check directory filters
			if !w.matchesFilters(path) {
				w.debug(path, "Directory filtered out", nil)
				return w.handleSkip(path, "filtered out")
			}

			if includeThis && w.Opts.IncludeStats && w.Stats != nil {
				w.mu.Lock()
				w.Stats.DirsVisited++
				w.Stats.Matches++
				w.mu.Unlock()
			}

			if includeThis && w.OnVisitDir != nil {
				if err := w.OnVisitDir(path, entries); err != nil {
					w.error(path, "OnVisitDir callback error", map[string]interface{}{"err": err})
					return w.handleError(path, err)
				}
			}

			return nil
		}

		// File handling - check filters BEFORE counting stats
		if !w.matchesFilters(path) {
			w.debug(path, "File filtered out", nil)
			return w.handleSkip(path, "filtered out")
		}

		// File handling
		info, err := d.Info()
		if err != nil {
			w.error(path, "Error getting file info", map[string]interface{}{"err": err})
			return w.handleError(path, err)
		}

		size := info.Size()
		w.trace(path, "File info retrieved", map[string]interface{}{
			"size":    size,
			"modTime": info.ModTime(),
		})

		// Update stats for matched files only
		if includeThis && w.Opts.IncludeStats && w.Stats != nil {
			w.mu.Lock()
			w.Stats.FilesVisited++
			w.Stats.TotalSize += size
			w.Stats.Matches++

			// Track min/max file sizes
			if w.Stats.MinFileSize == -1 || size < w.Stats.MinFileSize {
				w.Stats.MinFileSize = size
			}
			if size > w.Stats.MaxFileSize {
				w.Stats.MaxFileSize = size
			}
			w.mu.Unlock()
		}

		if includeThis && w.OnVisitFile != nil {
			if err := w.OnVisitFile(path, size); err != nil {
				w.error(path, "OnVisitFile callback error", map[string]interface{}{"err": err})
				return w.handleError(path, err)
			}
		}

		return nil
	})

	if w.Opts.IncludeStats && w.Stats != nil {
		w.Stats.EndTime = time.Now()
		w.Stats.Duration = w.Stats.EndTime.Sub(w.Stats.StartTime)

		if w.Stats.FilesVisited > 0 {
			w.Stats.AvgFileSize = w.Stats.TotalSize / int64(w.Stats.FilesVisited)
		}

		if w.Stats.Duration > 0 {
			w.Stats.EntriesPerSec = float64(w.Stats.EntriesVisited) / w.Stats.Duration.Seconds()
			w.Stats.FilesPerSec = float64(w.Stats.FilesVisited) / w.Stats.Duration.Seconds()
			w.Stats.DataRate = float64(w.Stats.TotalSize) / w.Stats.Duration.Seconds()
		}
	}

	w.debug(root, "Walk finished", map[string]interface{}{"error": err})

	return err
}

func (w *Walker) Count(root string) (*WalkStats, error) {
	// Save callbacks so we can restore after counting
	savedFileCb, savedDirCb := w.OnVisitFile, w.OnVisitDir
	savedProgress := w.Opts.EnableProgress

	w.OnVisitFile, w.OnVisitDir = nil, nil
	w.Opts.EnableProgress = false

	// Reset stats
	stats := &WalkStats{Custom: make(map[string]interface{})}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// ignore error, continue walking
			return nil
		}

		// Respect context cancellation
		select {
		case <-w.Context.Done():
			return w.Context.Err()
		default:
		}

		depth := getDepth(root, path)
		if w.Opts.SkipRoot && depth == 0 {
			return nil
		}

		// Skip entries shallower than MinIncludeDepth
		if w.Opts.MinIncludeDepth > 0 && depth < w.Opts.MinIncludeDepth {
			// continue walking children
			// if d.IsDir() {
			// 	return nil
			// }

			return nil
		}

		// Skip entries deeper than MaxDepth
		if w.Opts.MaxDepth >= 0 && depth > w.Opts.MaxDepth {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		// Skip hidden files/dirs unless IncludeHidden
		if !w.Opts.IncludeHidden && isHidden(d.Name()) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		// Apply filters (patterns, extensions, etc.)
		if !w.matchesFilters(path) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		// Skip empty dirs if configured
		// Skip empty dirs only if necessary
		if d.IsDir() && w.Opts.SkipEmptyDirs {
			entries, _ := os.ReadDir(path) // optional: ignore errors

			if len(entries) == 0 {
				return fs.SkipDir
			}
		}
		// if d.IsDir() {
		// 	entries, err := os.ReadDir(path)
		// 	if err == nil && w.Opts.SkipEmptyDirs && len(entries) == 0 {
		// 		return fs.SkipDir
		// 	}
		// }

		// Skip non-leaf dirs if OnlyLeafDirs
		if d.IsDir() && w.Opts.OnlyLeafDirs && !isLeafDir(path) {
			return nil
		}

		// Skip entries based on OnlyDirs / OnlyFiles
		if d.IsDir() && w.Opts.OnlyFiles {
			return nil // skip dir entirely
		}
		if !d.IsDir() && w.Opts.OnlyDirs {
			return nil // skip file entirely
		}

		// Passed all filters, increment counts
		stats.EntriesVisited++
		if d.IsDir() {
			stats.DirsVisited++
		} else {
			stats.FilesVisited++
		}

		return nil
	})

	// Restore callbacks and options
	w.OnVisitFile, w.OnVisitDir = savedFileCb, savedDirCb
	w.Opts.EnableProgress = savedProgress

	return stats, err
}

func (w *Walker) GetStats() *WalkStats {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Stats == nil {
		return &WalkStats{}
	}

	stats := *w.Stats

	// Ensure EndTime is set for accurate Duration
	if stats.EndTime.IsZero() {
		stats.EndTime = time.Now()
	}
	stats.Duration = stats.EndTime.Sub(stats.StartTime)

	// Compute averages if possible
	if stats.FilesVisited > 0 {
		stats.AvgFileSize = stats.TotalSize / int64(stats.FilesVisited)
	}

	// MinFileSize: default to 0 only if no files, otherwise set properly
	if stats.FilesVisited == 0 {
		stats.MinFileSize = 0
	}

	// Compute rates
	secs := stats.Duration.Seconds()
	if secs > 0 {
		stats.DataRate = float64(stats.TotalSize) / secs
		stats.EntriesPerSec = float64(stats.EntriesVisited) / secs
		stats.FilesPerSec = float64(stats.FilesVisited) / secs
	}

	return &stats
}

// matchesFilters checks include/exclude rules.
func (w *Walker) matchesFilters(path string) bool {
	// ext := strings.ToLower(filepath.Ext(path))
	// if len(w.Opts.IncludeExts) > 0 && !contains(w.Opts.IncludeExts, ext) {
	// 	fmt.Println("included exts")
	// 	return false
	// }

	// if len(w.Opts.ExcludeExts) > 0 && contains(w.Opts.ExcludeExts, ext) {
	// 	return false
	// }

	// Include pattern check
	if len(w.Opts.IncludePatterns) > 0 {
		matched := false
		for _, pat := range w.Opts.IncludePatterns {
			if ok, _ := filepath.Match(pat, filepath.Base(path)); ok {
				matched = true
				break
			}
		}

		if !matched {
			return false
		}
	}

	// Exclude pattern check
	if len(w.Opts.ExcludePatterns) > 0 {
		for _, pat := range w.Opts.ExcludePatterns {
			if ok, _ := filepath.Match(pat, filepath.Base(path)); ok {
				return false
			}
		}
	}

	return true
}

// handleError handles errors according to opts.
func (w *Walker) handleError(path string, err error) error {
	if w.Opts.IncludeStats && w.Stats != nil {
		w.mu.Lock()
		w.Stats.ErrorsCount++

		if w.Opts.IncludeErrors {
			w.Stats.Errors = append(w.Stats.Errors, err)
		}

		w.mu.Unlock()
	}

	if w.OnError != nil {
		if cbErr := w.OnError(path, err); cbErr != nil {
			return cbErr
		}
	}

	if w.Opts.StopOnError {
		return err
	}

	if w.Opts.SkipOnError {
		return nil
	}

	return err
}

// skip calls OnSkip if set.
func (w *Walker) handleSkip(path string, reason string) error {
	if w.Stats != nil {
		w.Stats.Skipped++
	}

	if w.OnSkip != nil {
		return w.OnSkip(path, reason)
	}

	return nil
}

func (w *Walker) injectPath(path string, detail any) map[string]interface{} {
	if detail == nil {
		return map[string]interface{}{"path": path}
	}

	if dm, ok := detail.(map[string]interface{}); ok {
		if _, exists := dm["path"]; !exists {
			dm["path"] = path
		}

		return dm
	}

	return map[string]interface{}{"path": path, "detail": detail}
}

// Log methods automatically inject path into detail
func (w *Walker) debug(path string, msg string, detail any) {
	if w.Logger == nil {
		return
	}

	w.Logger.Debug(msg, w.injectPath(path, detail))
}

func (w *Walker) info(path string, msg string, detail any) {
	if w.Logger == nil {
		return
	}

	w.Logger.Info(msg, w.injectPath(path, detail))
}

func (w *Walker) warn(path string, msg string, detail any) {
	if w.Logger == nil {
		return
	}

	w.Logger.Warn(msg, w.injectPath(path, detail))
}

func (w *Walker) error(path string, msg string, detail any) {
	if w.Logger == nil {
		return
	}

	w.Logger.Error(msg, w.injectPath(path, detail))
}

func (w *Walker) trace(path string, msg string, detail any) {
	if w.Logger == nil {
		return
	}

	w.Logger.Trace(msg, w.injectPath(path, detail))
}

func (s *WalkStats) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("=== Walk Stats ===\n")
	sb.WriteString(fmt.Sprintf("Dirs Visited:   %d\n", s.DirsVisited))
	sb.WriteString(fmt.Sprintf("Files Visited:  %d\n", s.FilesVisited))
	sb.WriteString(fmt.Sprintf("Errors Count:   %d\n", s.ErrorsCount))
	sb.WriteString(fmt.Sprintf("Total Size:     %d bytes\n", s.TotalSize))

	if len(s.Errors) > 0 {
		sb.WriteString("\n--- Errors ---\n")
		for i, err := range s.Errors {
			sb.WriteString(fmt.Sprintf("[%d] %v\n", i+1, err))
		}
	}

	return sb.String()
}

func (o WalkOptions) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("=== Walk Options ===\n")
	sb.WriteString(fmt.Sprintf("IncludeStats:   %t\n", o.IncludeStats))
	sb.WriteString(fmt.Sprintf("IncludeErrors:  %t\n", o.IncludeErrors))
	sb.WriteString(fmt.Sprintf("StopOnError:    %t\n", o.StopOnError))
	sb.WriteString(fmt.Sprintf("SkipOnError:    %t\n", o.SkipOnError))
	sb.WriteString(fmt.Sprintf("MaxDepth:       %d\n", o.MaxDepth))
	sb.WriteString(fmt.Sprintf("SkipEmpty:      %t\n", o.SkipEmptyDirs))
	// sb.WriteString(fmt.Sprintf("FollowSymlinks: %t\n", o.FollowSymlinks))
	sb.WriteString(fmt.Sprintf("IncludeHidden:  %t\n", o.IncludeHidden))

	return sb.String()
}

func FormatStats(s *WalkStats) string {
	// Format times nicely, or show "N/A" if zero
	start := "N/A"
	if !s.StartTime.IsZero() {
		start = s.StartTime.Format("2006-01-02 15:04:05")
	}

	end := "N/A"
	if !s.EndTime.IsZero() {
		end = s.EndTime.Format("2006-01-02 15:04:05")
	}

	// Calculate duration if zero but start and end are valid
	duration := s.Duration
	if duration == 0 && !s.StartTime.IsZero() && !s.EndTime.IsZero() {
		duration = s.EndTime.Sub(s.StartTime)
	}

	// Format Errors count and show first few errors if any
	errCount := len(s.Errors)
	var errSummary string
	if errCount == 0 {
		errSummary = "None"
	} else if errCount <= 3 {
		errSummary = ""
		for i, e := range s.Errors {
			errSummary += fmt.Sprintf("\n    %d: %v", i+1, e)
		}
	} else {
		errSummary = fmt.Sprintf("%d errors (showing first 3):", errCount)
		for i := 0; i < 3; i++ {
			errSummary += fmt.Sprintf("\n    %d: %v", i+1, s.Errors[i])
		}
	}

	// Format Custom map keys and values, or skip if empty
	customSummary := "None"
	if len(s.Custom) > 0 {
		customSummary = ""
		for k, v := range s.Custom {
			customSummary += fmt.Sprintf("\n    %s: %v", k, v)
		}
	}

	return fmt.Sprintf(`Scan Stats:
  Start Time:      %s
  End Time:        %s
  Duration:        %s

  Entries Visited: %d (Files: %d, Dirs: %d)
  Matches:         %d
  Skipped:         %d
  Errors Count:    %d
  Errors:          %s

  Total Size:      %d bytes
  Avg File Size:   %d bytes
  Min File Size:   %d bytes
  Max File Size:   %d bytes
  Data Rate:       %.2f bytes/sec

  Entries/sec:     %.2f
  Files/sec:       %.2f

  Custom Metrics:  %s
`, start, end, duration,
		s.EntriesVisited, s.FilesVisited, s.DirsVisited,
		s.Matches, s.Skipped, s.ErrorsCount, errSummary,
		s.TotalSize, s.AvgFileSize, s.MinFileSize, s.MaxFileSize, s.DataRate,
		s.EntriesPerSec, s.FilesPerSec,
		customSummary,
	)
}

func parseLevel(levelStr string) logger.Level {
	switch strings.ToUpper(levelStr) {
	case "ERROR":
		return logger.LevelError
	case "WARN":
		return logger.LevelWarn
	case "INFO":
		return logger.LevelInfo
	case "DEBUG":
		return logger.LevelDebug
	case "TRACE":
		return logger.LevelTrace
	default:
		return logger.LevelInfo
	}
}
