package scanner

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

// TODO: use scan options
func Scan[T any](
	sources map[string]string,
	cache *dirCache,
	opts ScanOptions,
	itemBuilder func(ScannedItem) (T, bool),
) []T {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []T
		tasks   = make(chan ScannedItem)
	)

	// === Progress Bar Setup ===
	var total int
	if opts.ShowProgress {
		for label, root := range sources {
			if opts.Mode == ScanDirs {
				total += CountTargetDirs(root, label, opts.Depth, cache, func(path string, entries []os.DirEntry) bool {
					return len(entries) > 0 // Skip empty dirs
				})
			} else if opts.Mode == ScanFiles {
				// Optional: count matching files here if needed
			}
		}
	}
	progress := NewProgress(total, opts.ShowProgress, "Scanning")

	// Worker pool
	for i := 0; i < opts.Concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for item := range tasks {
				if out, ok := itemBuilder(item); ok {
					mu.Lock()
					results = append(results, out)
					mu.Unlock()
				}

				progress.Add(1)
			}
		}()
	}

	for source, root := range sources {
		switch opts.Mode {
		case ScanDirs:
			_ = walkDirs(root, opts.Depth, cache, func(dirPath string, entries []os.DirEntry) {
				groupLabel := buildGroupLabel(root, dirPath)

				tasks <- ScannedItem{
					Source:     source,
					GroupLabel: groupLabel,
					GroupPath:  root,
					ItemPath:   dirPath,
					ItemName:   filepath.Base(dirPath),
					SubEntries: entries,
				}
			})

		case ScanFiles:
			_ = walkFiles(root, opts.Depth, opts.Exts, cache, func(filePath string) {
				groupLabel := buildGroupLabel(root, filepath.Dir(filePath))

				tasks <- ScannedItem{
					Source:     source,
					GroupLabel: groupLabel,
					GroupPath:  root,
					ItemPath:   filePath,
					ItemName:   filepath.Base(filePath),
					SubEntries: nil,
				}
			})

		default:
			log.Fatalf("Scan(): unsupported scan mode %q for source %s", opts.Mode, source)

		}
	}

	close(tasks)
	wg.Wait()
	progress.Finish()

	return results

	// // Count all group directories first for progress bar
	// var totalGroups int
	// for _, path := range sources {
	// 	if entries := cache.Read(path); entries != nil {
	// 		for _, entry := range entries {
	// 			if entry.IsDir() {
	// 				totalGroups++
	// 			}
	// 		}
	// 	}
	// }

	// bar := progressbar.NewOptions(totalGroups,
	// 	progressbar.OptionSetDescription("Scanning"),
	// 	progressbar.OptionShowCount(),
	// 	progressbar.OptionSetWidth(20),
	// 	progressbar.OptionClearOnFinish(),
	// )

	// for label, path := range sources {
	// 	groupDirs := cache.Read(path)
	// 	if groupDirs == nil {
	// 		continue
	// 	}

	// 	for _, group := range groupDirs {
	// 		if !group.IsDir() {
	// 			continue
	// 		}

	// 		groupPath := filepath.Join(path, group.Name())

	// 		wg.Add(1)
	// 		go func(groupName, groupPath string) {
	// 			defer wg.Done()
	// 			sem <- struct{}{}        // 🛑 acquire
	// 			defer func() { <-sem }() // ✅ release

	// 			dirEntries := cache.Read(groupPath)
	// 			if dirEntries == nil {
	// 				_ = bar.Add(1)
	// 				return
	// 			}

	// 			for _, entry := range dirEntries {
	// 				if !entry.IsDir() {
	// 					continue
	// 				}

	// 				itemPath := filepath.Join(groupPath, entry.Name())
	// 				subEntries := cache.Read(itemPath)
	// 				if subEntries == nil {
	// 					continue
	// 				}

	// 				if item, ok := itemBuilder(itemPath, label, subEntries); ok {
	// 					mu.Lock()
	// 					results = append(results, item)
	// 					mu.Unlock()
	// 				}
	// 			}

	// 			_ = bar.Add(1)

	// 		}(group.Name(), groupPath)
	// 	}
	// }

	// wg.Wait()
	// bar.Finish()

	// sort.Slice(results, func(i, j int) bool {
	// 	return results[i].GetName() < results[j].GetName()
	// })

	// return results
}
