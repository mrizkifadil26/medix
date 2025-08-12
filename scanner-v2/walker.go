package scannerV2

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Walker struct {
	Context  context.Context
	Opts     WalkOptions
	Stats    *WalkStats
	Progress *ProgressTracker

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
	Enable  bool
	LogFunc func(event DebugEvent)
}

// DebugEvent is a structured debug message from the walker.
type DebugEvent struct {
	Time    time.Time
	Level   string // always "DEBUG" here
	Path    string // affected file/dir path
	Message string // short description
	Detail  any    // optional extra info (filters, depth, etc.)
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

func NewWalker(ctx context.Context, opts WalkOptions) *Walker {
	return &Walker{
		Context: ctx,
		Opts:    opts,
		Stats:   &WalkStats{Custom: make(map[string]interface{})},
	}
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
		stats, _ := w.Count(root)
		total := stats.EntriesVisited
		w.debug(root, "Total entries found", map[string]interface{}{"total": total})

		if total > 0 {
			w.Progress = NewProgressTracker(total, true, "Scanning")
			w.debug(root, "Finishing progress", nil)

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

			w.debug(root, "Stats collected", w.Stats)
		}
	}()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
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
			w.debug(path, "Skipping hidden entry", nil)
			return w.handleSkip(path, "hidden entry")
		}

		if w.Opts.SkipRoot && depth == 0 {
			// Root: skip processing, but allow traversal into immediate children
			return nil
		}

		if w.Opts.MaxDepth >= 0 && depth > w.Opts.MaxDepth {
			w.debug(path, "Skipping due to depth limit", nil)
			return fs.SkipDir
		}

		// Skip based on OnlyDirs / OnlyFiles flags:
		if w.Opts.OnlyDirs && !d.IsDir() {
			return nil // skip files
		}
		if w.Opts.OnlyFiles && d.IsDir() {
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
				w.debug(path, "Error reading directory", map[string]interface{}{"err": err})
				return w.handleError(path, err)
			}

			// Skip empty dirs
			if w.Opts.SkipEmptyDirs && len(entries) == 0 {
				w.debug(path, "Skipping empty dir", nil)
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
					w.debug(path, "OnVisitDir returned error", map[string]interface{}{"err": err})
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
			w.debug(path, "Error getting file info", map[string]interface{}{"err": err})
			return w.handleError(path, err)
		}

		size := info.Size()

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
				w.debug(path, "OnVisitFile returned error", map[string]interface{}{"err": err})
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

	return err
}

func (w *Walker) Count(root string) (*WalkStats, error) {
	// Save callbacks so we can restore after counting
	savedFileCb := w.OnVisitFile
	savedDirCb := w.OnVisitDir
	savedProgress := w.Opts.EnableProgress

	w.OnVisitFile = nil
	w.OnVisitDir = nil
	w.Opts.EnableProgress = false

	// Reset stats
	w.Stats = &WalkStats{Custom: make(map[string]interface{})}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip error during counting
		}

		// Respect context cancellation
		select {
		case <-w.Context.Done():
			return w.Context.Err()
		default:
		}

		// Skip hidden
		if !w.Opts.IncludeHidden && isHidden(d.Name()) {
			return nil
		}

		// Skip depth beyond MaxDepth
		depth := getDepth(root, path)
		// if depth >= w.Opts.MinIncludeDepth {
		// 	w.Stats.EntriesVisited++
		// }

		if w.Opts.SkipRoot && depth == 0 {
			// Root: skip processing, but allow traversal into immediate children
			return nil
		}

		// Enforce min include depth if set
		if w.Opts.MinIncludeDepth > 0 && depth < w.Opts.MinIncludeDepth {
			// We don't count entries shallower than min depth
			// but continue traversal
			return nil
		}

		if w.Opts.MaxDepth >= 0 && depth > w.Opts.MaxDepth {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		// Skip hidden entries if needed
		if !w.Opts.IncludeHidden && isHidden(d.Name()) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		// Filter by patterns and extensions
		if !w.matchesFilters(path) {
			if d.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if d.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil // ignore error in count mode
			}

			// Skip empty dirs
			if w.Opts.SkipEmptyDirs && len(entries) == 0 {
				return nil
			}

			// Only leaf dirs
			if w.Opts.OnlyLeafDirs && d.IsDir() && !isLeafDir(path) {
				return nil
			}

			// if !w.matchesFilters(path) {
			// 	return nil
			// }

			w.Stats.DirsVisited++
			w.Stats.EntriesVisited++
			return nil
		}

		// File filtering
		// if !w.matchesFilters(path) {
		// 	return nil
		// }

		w.Stats.FilesVisited++
		w.Stats.EntriesVisited++
		return nil
	})

	// Restore callbacks
	w.OnVisitFile = savedFileCb
	w.OnVisitDir = savedDirCb
	w.Opts.EnableProgress = savedProgress

	return w.Stats, err
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
	ext := strings.ToLower(filepath.Ext(path))
	if len(w.Opts.IncludeExts) > 0 && !contains(w.Opts.IncludeExts, ext) {
		return false
	}

	if len(w.Opts.ExcludeExts) > 0 && contains(w.Opts.ExcludeExts, ext) {
		return false
	}

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

// debug emits debug events if enabled.
func (w *Walker) debug(path, msg string, detail any) {
	if w.Opts.Debug.Enable && w.Opts.Debug.LogFunc != nil {
		w.Opts.Debug.LogFunc(DebugEvent{
			Time:    time.Now(),
			Level:   "DEBUG",
			Path:    path,
			Message: msg,
			Detail:  detail,
		})
	}
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
