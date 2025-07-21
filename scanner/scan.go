package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/mrizkifadil26/medix/model"
	"github.com/schollz/progressbar/v3"
)

func scanMedia(
	sources map[string]string,
	cache *dirCache,
	itemBuilder func(titlePath, label string, subEntries []os.DirEntry) (model.MediaEntry, bool),
	concurrency int,
) []model.MediaEntry {
	var (
		mu      sync.Mutex
		results []model.MediaEntry
		wg      sync.WaitGroup
		sem     = make(chan struct{}, concurrency)
	)

	// Count all group directories first for progress bar
	var totalGroups int
	for _, path := range sources {
		if entries := cache.Read(path); entries != nil {
			for _, entry := range entries {
				if entry.IsDir() {
					totalGroups++
				}
			}
		}
	}

	bar := progressbar.NewOptions(totalGroups,
		progressbar.OptionSetDescription("Scanning"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(20),
		progressbar.OptionClearOnFinish(),
	)

	for label, path := range sources {
		groupDirs := cache.Read(path)
		if groupDirs == nil {
			continue
		}

		for _, group := range groupDirs {
			if !group.IsDir() {
				continue
			}

			groupPath := filepath.Join(path, group.Name())

			wg.Add(1)
			go func(groupName, groupPath string) {
				defer wg.Done()
				sem <- struct{}{}        // 🛑 acquire
				defer func() { <-sem }() // ✅ release

				dirEntries := cache.Read(groupPath)
				if dirEntries == nil {
					_ = bar.Add(1)
					return
				}

				for _, entry := range dirEntries {
					if !entry.IsDir() {
						continue
					}

					itemPath := filepath.Join(groupPath, entry.Name())
					subEntries := cache.Read(itemPath)
					if subEntries == nil {
						continue
					}

					if item, ok := itemBuilder(itemPath, label, subEntries); ok {
						mu.Lock()
						results = append(results, item)
						mu.Unlock()
					}
				}

				_ = bar.Add(1)

			}(group.Name(), groupPath)
		}
	}

	wg.Wait()
	bar.Finish()

	sort.Slice(results, func(i, j int) bool {
		return results[i].GetName() < results[j].GetName()
	})

	return results
}
