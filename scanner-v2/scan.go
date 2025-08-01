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
		jobs  []concurrency.TaskFunc
		mu    sync.Mutex // to protect shared output
		items []ScanEntry

		excluded int64 // atomic counter
	)

	switch options.Mode {
	case "files":
		err = WalkFiles(inputPath, WalkOptions{
			MaxDepth: options.Depth,
			Exts:     options.Exts,
			Verbose:  options.Verbose,
		}, func(path string, size int64) {
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
		})

	case "dirs":
		err = WalkDirs(
			inputPath,
			WalkOptions{
				MaxDepth:  options.Depth,
				OnlyLeaf:  options.OnlyLeaf,
				LeafDepth: options.LeafDepth,
				SkipEmpty: options.SkipEmpty,
				Verbose:   options.Verbose,
			},
			func(path string, entries []os.DirEntry) {
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
							if !options.SubEntries {
								return nil
							}

							var subs []ScanEntry
							extFilter := make(map[string]bool)
							for _, ext := range options.SubExts {
								extFilter[strings.ToLower(ext)] = true
							}

							for _, e := range entries {
								if e.IsDir() {
									continue
								}

								name := e.Name()
								ext := strings.ToLower(filepath.Ext(name))
								if len(extFilter) > 0 && !extFilter[ext] {
									continue
								}

								info, err := e.Info()
								if err != nil {
									continue // safely skip files we can't stat
								}
								size := info.Size()

								subs = append(subs, ScanEntry{
									ItemPath: filepath.Join(path, e.Name()),
									ItemName: e.Name(),
									ItemSize: &size,
								})
							}

							if len(subs) > 0 {
								return subs
							}

							return nil
						}(),
					}

					mu.Lock()
					items = append(items, entry)
					mu.Unlock()

					return nil
				})
			})

	default:
		return output, fmt.Errorf("unsupported scan mode: %s", options.Mode)
	}

	if err != nil {
		return output, fmt.Errorf("scan failed: %w", err)
	}

	taskExec, err := SelectExecutor(options.Concurrency)
	if err != nil {
		return output, fmt.Errorf("concurrency error: %w", err)
	}

	exec := concurrency.FromTaskExecutor(taskExec)
	err = exec(ctx, jobs)
	if err != nil {
		return output, fmt.Errorf("execution error: %w", err)
	}

	output.Items = items
	output.ItemCount = len(items)
	output.ExcludedCount = int(atomic.LoadInt64(&excluded))
	output.Duration = time.Since(start).String()

	return output, nil
}

func (o ScanOptions) IsParallel() bool {
	return o.Concurrency > 1
}
