package scannerV2

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Walker struct {
	Context context.Context
	Opts    WalkOptions
	Stats   *WalkStats

	OnVisitFile func(path string, size int64) error
	OnVisitDir  func(path string, entries []fs.DirEntry) error
	OnSkip      func(path string, reason string) error
	OnError     func(path string, err error) error
}

type WalkOptions struct {
	StopOnError bool
	SkipOnError bool

	MaxDepth      int
	SkipEmptyDirs bool
	OnlyLeafDirs  bool
	IncludeHidden bool

	IncludePatterns []string
	ExcludePatterns []string
	IncludeExts     []string
	ExcludeExts     []string

	EnableStats bool // collect walk statistics

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

func New(ctx context.Context, opts WalkOptions) *Walker {
	return &Walker{
		Context: ctx,
		Opts:    opts,
		Stats:   &WalkStats{Custom: make(map[string]interface{})},
	}
}

func (w *Walker) Walk(root string) error {
	start := time.Now()

	// Initialize stats if enabled
	if w.Opts.EnableStats && w.Stats != nil {
		w.Stats.StartTime = start
		w.Stats.MinFileSize = -1 // sentinel for tracking min
		if w.Stats.Custom == nil {
			w.Stats.Custom = make(map[string]interface{})
		}
	}

	defer func() {
		if w.Opts.EnableStats && w.Stats != nil {
			w.Stats.EndTime = time.Now()
			w.Stats.Duration = w.Stats.EndTime.Sub(w.Stats.StartTime)
			if w.Stats.FilesVisited > 0 {
				w.Stats.AvgFileSize = w.Stats.TotalSize / int64(w.Stats.FilesVisited)
			}

			if w.Stats.Duration > 0 {
				durSec := w.Stats.Duration.Seconds()
				w.Stats.EntriesPerSec = float64(w.Stats.EntriesVisited) / durSec
				w.Stats.FilesPerSec = float64(w.Stats.FilesVisited) / durSec
				w.Stats.DataRate = float64(w.Stats.TotalSize) / durSec
			}
		}
	}()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return w.handleError(path, err)
		}

		// Respect context cancellation
		select {
		case <-w.Context.Done():
			w.handleSkip(path, "context canceled")
			return w.handleError(path, w.Context.Err())
		default:
		}

		// --- Stats: entries visited ---
		if w.Opts.EnableStats && w.Stats != nil {
			w.Stats.EntriesVisited++
		}

		// Hidden file/dir skipping
		if !w.Opts.IncludeHidden && isHidden(d.Name()) {
			return w.handleSkip(path, "hidden entry")
		}

		depth := getDepth(root, path)
		if w.Opts.MaxDepth >= 0 && depth > w.Opts.MaxDepth {
			w.debug(path, "Skipping due to depth limit", nil)
			return fs.SkipDir
		}

		// Directory handling
		if d.IsDir() {
			if w.Opts.EnableStats {
				w.Stats.DirsVisited++
			}

			entries, err := os.ReadDir(path)
			if err != nil {
				return w.handleError(path, err)
			}

			// Skip empty dirs
			if w.Opts.SkipEmptyDirs && len(entries) == 0 {
				return w.handleSkip(path, "empty dir")
			}

			// Only leaf dirs (skip if not leaf)
			if w.Opts.OnlyLeafDirs && len(entries) > 0 {
				return nil
			}

			if w.OnVisitDir != nil {
				if err := w.OnVisitDir(path, entries); err != nil {
					return err
				}
			}

			return nil
		}

		// File handling
		info, err := d.Info()
		if err != nil {
			return w.handleError(path, err)
		}

		size := info.Size()

		w.Stats.FilesVisited++
		w.Stats.TotalSize += size

		if !w.matchesFilters(path) {
			return w.handleSkip(path, "filtered out")
		}

		w.Stats.Matches++
		if w.OnVisitFile != nil {
			if err := w.OnVisitFile(path, size); err != nil {
				return err
			}
		}

		return nil
	})

	if w.Opts.EnableStats && w.Stats != nil {
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
	if w.Opts.EnableStats && w.Stats != nil {
		w.Stats.ErrorsCount++
		w.Stats.Errors = append(w.Stats.Errors, err)
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
	w.Stats.Skipped++
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
