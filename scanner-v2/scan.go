package scannerV2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mrizkifadil26/medix/utils/concurrency"
)

func Scan(root string, options ScanOptions) (ScanOutput, error) {
	ctx := context.Background()

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

	var jobs []concurrency.TaskFunc
	var mu sync.Mutex // to protect shared output
	excludedCount := 0

	switch options.Mode {
	case "files":
		err = WalkFiles(inputPath, WalkOptions{
			MaxDepth: options.Depth,
			Exts:     options.Exts,
			Verbose:  options.Verbose,
		}, func(path string, size int64) {
			if isExcluded(path) {
				excludedCount++
				return
			}

			rel, _ := filepath.Rel(inputPath, path)
			jobs = append(jobs, concurrency.TaskFunc(func(ctx context.Context) error {
				entry := ScanEntry{
					// Source:     inputPath,
					GroupPath:  filepath.Dir(rel),
					ItemPath:   path,
					ItemName:   filepath.Base(path),
					ItemSize:   &size,
					GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
				}

				mu.Lock()
				output.Items = append(output.Items, entry)
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
				excludedCount := 0
				if isExcluded(path) {
					excludedCount++
					return
				}

				rel, _ := filepath.Rel(inputPath, path)

				jobs = append(jobs, func(ctx context.Context) error {
					entry := ScanEntry{
						GroupPath:  filepath.Dir(rel),
						ItemPath:   path,
						ItemName:   filepath.Base(path),
						GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
						SubEntries: func() []string {
							var subs []string
							for _, e := range entries {
								if e.IsDir() {
									subs = append(subs, e.Name())
								}
							}
							if len(subs) > 0 {
								return subs
							}
							return nil
						}(),
					}

					mu.Lock()
					output.Items = append(output.Items, entry)
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

	// Use executor
	taskExec, err := SelectExecutor(options.Concurrency)
	if err != nil {
		return output, fmt.Errorf("concurrency error: %w", err)
	}

	exec := concurrency.FromTaskExecutor(taskExec)
	err = exec(ctx, jobs)
	if err != nil {
		return output, fmt.Errorf("execution error: %w", err)
	}

	output.ExcludedCount = excludedCount
	output.ItemCount = len(output.Items)
	output.Duration = time.Since(start).String()

	return output, nil
}

func (o ScanOptions) IsParallel() bool {
	return o.Concurrency > 1
}
