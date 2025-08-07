package scannerV2

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Walker struct {
	Ctx   context.Context
	Root  string
	Opts  WalkOptions
	Stats *WalkerStats // Collects walk statistics

	// File matching / filtering
	IsExcluded func(path string) bool
	MatchExt   func(name string) bool
	MatchDir   func(path string, entry fs.DirEntry) bool

	// Callbacks
	OnVisitFile func(path string, size int64)
	OnVisitDir  func(path string, entries []fs.DirEntry)
	OnSkip      func(path string, reason string)

	Logger func(format string, args ...any) // Optional custom logger (defaults to log.Printf)

	// Caching / State
	Cache map[string]fs.FileInfo // Optional metadata cache (e.g., size, modtime)
	mutex *sync.Mutex            // Guards concurrent access to shared fields
	trace []string               // Optional trace of visited paths (for debugging or test logs)
}

type WalkOptions struct {
	// Control behavior
	MaxDepth int

	// Filtering
	Exts      []string // match file extensions
	OnlyLeaf  bool     // only include leaf directories
	LeafDepth int      // 0 = default, 1 = leaf-1, etc.
	SkipEmpty bool     // skip empty directories

	// Logging control
	Verbose   bool // high-level logs
	DebugMode bool // low-level/internal logs
}

type WalkerStats struct {
	VisitedFiles int64
	VisitedDirs  int64
	Skipped      int64
	Excluded     int64
	MatchedFiles int64
	EmptyDirs    int64

	CacheHits   int64
	CacheMisses int64

	StartTime time.Time
	EndTime   time.Time
}

func WalkDirs(root string, opts WalkOptions, fn func(path string, entries []os.DirEntry)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Error accessing %s: %v\n", path, err)
			}

			return err
		}
		if !d.IsDir() {
			return nil
		}

		depth := pathDepth(root, path)
		if opts.MaxDepth >= 0 && depth > opts.MaxDepth {
			if opts.Verbose {
				fmt.Printf("⏭️  Skipping %s (too deep, depth %d)\n", path, depth)
			}

			return fs.SkipDir
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Failed to read dir %s: %v\n", path, err)
			}

			return err
		}

		// Skip empty if requested
		if opts.SkipEmpty && len(entries) == 0 {
			if opts.Verbose {
				fmt.Printf("⏭️  Skipping %s (empty directory)\n", path)
			}

			return nil
		}

		// Handle leaf-depth scan mode
		if opts.LeafDepth > 0 {
			leafLevel, err := getLeafDepth(path)
			if err != nil {
				if opts.Verbose {
					fmt.Printf("⚠️  Failed to evaluate leaf depth of %s: %v\n", path, err)
				}
				return err
			}
			if leafLevel != opts.LeafDepth {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping %s (leafDepth=%d != expected=%d)\n", path, leafLevel, opts.LeafDepth)
				}
				return nil
			}

			// Special case: skip group-level node at exact depth-1
			if opts.LeafDepth == 1 && path == root {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping root %s for leafDepth=1\n", path)
				}

				return nil
			}
		} else if opts.OnlyLeaf {
			if containsDir(entries) {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping %s (not a leaf folder)\n", path)
				}
				return nil
			}

			// If LeafDepth == 1 + OnlyLeaf: only visit depth 1 and skip traversal
			if opts.MaxDepth == 1 && depth == 1 {
				fn(path, entries)
				return fs.SkipDir
			}
		}

		if opts.Verbose {
			fmt.Printf("📂 Visiting %s (%d entries)\n", path, len(entries))
		}

		fn(path, entries)
		return nil
	})
}

func WalkFiles(root string, opts WalkOptions, fn func(path string, size int64)) error {
	extMap := make(map[string]bool)
	for _, ext := range opts.Exts {
		extMap[strings.ToLower(ext)] = true
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Error accessing %s: %v\n", path, err)
			}
			return err
		}

		if d.IsDir() {
			depth := pathDepth(root, path)
			if opts.MaxDepth >= 0 && depth > opts.MaxDepth {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping dir %s (too deep, depth %d)\n", path, depth)
				}

				return fs.SkipDir
			}

			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if len(extMap) > 0 && !extMap[ext] {
			if opts.Verbose {
				fmt.Printf("⏭️  Skipping file %s (unsupported ext: %s)\n", path, ext)
			}

			return nil
		}

		info, err := d.Info()
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Failed to get info for file %s: %v\n", path, err)
			}

			return nil
		}

		if opts.Verbose {
			fmt.Printf("📄 Found file %s (%d bytes)\n", path, info.Size())
		}

		fn(path, info.Size())
		return nil
	})
}

func (w *Walker) log(format string, args ...any) {
	if w.Opts.Verbose {
		if w.Logger != nil {
			w.Logger(format, args...)
		} else {
			log.Printf(format, args...)
		}
	}
}

func (w *Walker) trackSkip(path, reason string) {
	if w.Stats != nil {
		w.Stats.Skipped++
	}

	w.appendTrace(path)

	if w.OnSkip != nil {
		w.OnSkip(path, reason)
	}

	w.log("⏭️  Skipping %s (%s)", path, reason)
}

func (w *Walker) appendTrace(path string) {
	// if w.Trace != nil {
	// 	w.Trace = append(w.Trace, path)
	// }
}

// Duration returns total time spent walking.
func (ws *WalkerStats) Duration() time.Duration {
	if !ws.StartTime.IsZero() && !ws.EndTime.IsZero() {
		return ws.EndTime.Sub(ws.StartTime)
	}

	return 0
}

func (w *Walker) Walk() error {
	if w.Stats != nil {
		w.Stats.StartTime = time.Now()
		defer func() {
			w.Stats.EndTime = time.Now()
		}()
	}

	return filepath.WalkDir(w.Root, func(path string, d fs.DirEntry, err error) error {
		// TODO: handle walk errors gracefully
		if err != nil {
			w.log("⚠️ Error accessing %s: %v", path, err)
			w.trackSkip(path, "walk error")
			return err
		}

		// Respect context cancellation
		select {
		case <-w.Ctx.Done():
			w.trackSkip(path, "context canceled")
			return w.Ctx.Err()
		default:
		}

		// Apply exclude filter
		if w.IsExcluded != nil && w.IsExcluded(path) {
			if w.Stats != nil {
				w.Stats.Excluded++
			}

			w.trackSkip(path, "excluded")
			return fs.SkipDir
		}

		depth := pathDepth(w.Root, path)
		if w.Opts.MaxDepth >= 0 && depth > w.Opts.MaxDepth {
			w.trackSkip(path,
				fmt.Sprintf("depth %d > maxDepth %d", depth, w.Opts.MaxDepth),
			)

			return fs.SkipDir
		}

		var info fs.FileInfo
		if w.Cache != nil {
			w.mutex.Lock()
			info = w.Cache[path]
			w.mutex.Unlock()
		}
		if info == nil {
			var err error
			info, err = d.Info()
			if err != nil {
				w.trackSkip(path, "info error")
				return nil
			}

			if w.Cache != nil {
				w.mutex.Lock()
				w.Cache[path] = info
				w.mutex.Unlock()
			}
		}

		// File visit
		if !d.IsDir() {
			if w.Stats != nil {
				w.Stats.VisitedFiles++
			}

			if w.MatchExt != nil && !w.MatchExt(d.Name()) {
				w.trackSkip(path, "extension mismatch")
				return nil
			}

			info, err := d.Info()
			if err != nil {
				w.trackSkip(path, "info error")
				return nil
			}

			if w.Stats != nil {
				w.Stats.MatchedFiles++
			}

			if w.OnVisitFile != nil {
				w.OnVisitFile(path, info.Size())
			}

			return nil
		}

		// Directory visit
		if w.Stats != nil {
			w.Stats.VisitedDirs++
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			w.trackSkip(path, "readdir error")
			return err
		}

		if w.Opts.SkipEmpty && len(entries) == 0 {
			w.trackSkip(path, "empty")
			if w.Stats != nil {
				w.Stats.EmptyDirs++
			}
			return nil
		}

		// Leaf-level filter
		if w.Opts.LeafDepth > 0 {
			leafLevel, err := getLeafDepth(path)
			if err != nil {
				w.trackSkip(path, "leaf depth check failed")
				return err
			}

			if leafLevel != w.Opts.LeafDepth {
				w.trackSkip(path, fmt.Sprintf("leafDepth=%d != expected=%d", leafLevel, w.Opts.LeafDepth))
				return nil
			}

			if w.Opts.LeafDepth == 1 && path == w.Root {
				w.trackSkip(path, "skipping root for leafDepth=1")
				return nil
			}
		} else if w.Opts.OnlyLeaf {
			if containsDir(entries) {
				w.trackSkip(path, "not a leaf folder")
				return nil
			}
		}

		if w.OnVisitDir != nil {
			fmt.Printf("📂 Visiting %s (%d entries)\n", path, len(entries))
			w.OnVisitDir(path, entries)
		}

		return nil
	})
}

func (ws *WalkerStats) PrintSummary() {
	fmt.Printf("📊 Walk Summary:\n")
	fmt.Printf("  Files visited:      %d\n", ws.VisitedFiles)
	fmt.Printf("  Matched files:      %d\n", ws.MatchedFiles)
	fmt.Printf("  Directories visited:%d\n", ws.VisitedDirs)
	fmt.Printf("  Empty directories:  %d\n", ws.EmptyDirs)
	fmt.Printf("  Excluded:           %d\n", ws.Excluded)
	fmt.Printf("  Skipped:            %d\n", ws.Skipped)
	fmt.Printf("  Duration:           %s\n", ws.Duration())

	fmt.Println()
}
