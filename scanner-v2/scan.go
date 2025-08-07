package scannerV2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mrizkifadil26/medix/utils/concurrency"
)

func Scan(root string, options ScanOptions) (ScanOutput, error) {
	// ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	start := time.Now()
	output := ScanOutput{
		GeneratedAt: time.Now().Format(time.RFC3339),
		SourcePath:  root,
		Mode:        options.Mode,
	}

	// Normalize input path
	inputPath := filepath.Clean(root)
	info, err := os.Stat(inputPath)
	if err != nil {
		return output, fmt.Errorf("input path error: %w", err)
	}

	if !info.IsDir() {
		return output, fmt.Errorf("input path is not a directory")
	}

	// Normalize exclude list
	excludeMap := make(map[string]bool)
	for _, ex := range options.Exclude {
		excludeMap[filepath.ToSlash(filepath.Clean(ex))] = true
	}

	// Helper: check if path is excluded
	isExcluded := func(absPath string) bool {
		rel, err := filepath.Rel(root, absPath)
		if err != nil {
			return false // fallback: don’t exclude
		}

		rel = filepath.ToSlash(rel) // make sure slashes match
		for prefix := range excludeMap {
			if strings.HasPrefix(rel, prefix) {
				return true
			}
		}

		return false
	}

	var (
		items    = make([]ScanEntry, 0)
		jobs     []concurrency.TaskFunc
		mu       sync.Mutex // to protect shared output
		stats    WalkerStats
		excluded int64 // atomic counter
	)

	// Extension filter for files
	matchExt := func(name string) bool {
		if len(options.Exts) == 0 {
			return true
		}
		ext := strings.ToLower(filepath.Ext(name))
		for _, e := range options.Exts {
			if strings.ToLower(e) == ext {
				return true
			}
		}
		return false
	}

	walker := &Walker{
		Root:       inputPath,
		Ctx:        ctx,
		Stats:      &stats,
		mutex:      &mu,
		IsExcluded: isExcluded,
		MatchExt:   matchExt,
		Opts: WalkOptions{
			MaxDepth: options.Depth,
			Exts:     options.Exts,
			Verbose:  options.Verbose,
		},
	}

	switch options.Mode {
	case "files":
		walker.OnVisitFile = func(
			path string,
			size int64,
		) {
			if isExcluded(path) {
				atomic.AddInt64(&excluded, 1)
				return
			}

			rel, _ := filepath.Rel(inputPath, path)

			jobs = append(jobs, concurrency.TaskFunc(func(ctx context.Context) error {
				entry := ScanEntry{
					GroupPath:  filepath.Dir(rel),
					ItemPath:   path,
					ItemName:   filepath.Base(path),
					ItemSize:   &size,
					GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
				}

				mu.Lock()
				items = append(items, entry)
				mu.Unlock()

				return nil
			}))
		}

	case "dirs":
		walker.OnVisitDir = func(
			path string,
			entries []os.DirEntry,
		) {
			if isExcluded(path) {
				atomic.AddInt64(&excluded, 1)
				return
			}

			rel, _ := filepath.Rel(inputPath, path)

			jobs = append(jobs, func(ctx context.Context) error {
				entry := ScanEntry{
					GroupPath:  filepath.Dir(rel),
					ItemPath:   path,
					ItemName:   filepath.Base(path),
					GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
					SubEntries: func() []ScanEntry {
						switch options.SubEntries {
						case SubentriesNone:
							return nil
						case SubentriesFlat:
							return scanFlat(path, options)
						case SubentriesNested:
							return scanNested(path, options)
						case SubentriesAuto:
							// Auto mode: if subdepth is -1, use nested, otherwise flat
							return nil
						default:
							return nil
						}
					}(),
				}

				mu.Lock()
				items = append(items, entry)
				mu.Unlock()

				return nil
			})
		}

	default:
		return output, fmt.Errorf("unsupported scan mode: %s", options.Mode)
	}

	if err := walker.Walk(); err != nil {
		return output, fmt.Errorf("scan failed: %w", err)
	}

	// if err != nil {
	// 	return output, fmt.Errorf("scan failed: %w", err)
	// }

	taskExec, err := SelectExecutor(options.Concurrency)
	if err != nil {
		return output, fmt.Errorf("concurrency error: %w", err)
	}

	exec := concurrency.FromTaskExecutor(taskExec)
	err = exec(ctx, jobs)
	if err != nil {
		return output, fmt.Errorf("execution error: %w", err)
	}

	walker.Stats.PrintSummary()

	output.Items = items
	output.ItemCount = len(items)
	output.ExcludedCount = int(atomic.LoadInt64(&excluded))
	output.Duration = time.Since(start).String()
	output.WalkStatistics = walker.Stats

	return output, nil
}

func (o ScanOptions) IsParallel() bool {
	return o.Concurrency > 1
}

func scanFlat(base string, opts ScanOptions) []ScanEntry {
	entries, _ := os.ReadDir(base)
	var out []ScanEntry

	// Create lowercase extension filter map (same as collectSubEntries)
	extFilter := make(map[string]bool)
	for _, ext := range opts.SubExts {
		extFilter[strings.ToLower(ext)] = true
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if len(extFilter) > 0 && !extFilter[ext] {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil
		}
		size := info.Size()

		fullPath := filepath.Join(base, entry.Name())
		out = append(out, ScanEntry{
			ItemPath: fullPath,
			ItemName: entry.Name(),
			ItemSize: &size,
		})
	}

	return out
}

func scanNested(base string, opts ScanOptions) []ScanEntry {
	return scanRecursive(base, 0, opts)
}

func scanRecursive(path string, depth int, opts ScanOptions) []ScanEntry {
	if opts.SubDepth >= 0 && depth >= opts.SubDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	// Lowercase extension filter
	extFilter := make(map[string]bool)
	for _, ext := range opts.SubExts {
		extFilter[strings.ToLower(ext)] = true
	}

	var out []ScanEntry

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			sub := scanRecursive(full, depth+1, opts)
			if opts.SkipEmpty && len(sub) == 0 {
				continue
			}

			out = append(out, ScanEntry{
				ItemPath:   full,
				ItemName:   entry.Name(),
				SubEntries: sub,
			})
		} else {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if len(extFilter) > 0 && !extFilter[ext] {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue // skip if we can't get info
			}
			size := info.Size()

			out = append(out, ScanEntry{
				ItemPath: full,
				ItemName: entry.Name(),
				ItemSize: &size,
			})
		}
	}

	return out
}
