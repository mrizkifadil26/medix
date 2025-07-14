package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/mrizkifadil26/medix/model"
	"github.com/schollz/progressbar/v3"
)

func scanGenericGroup[T any, G any](
	roots []string,
	cache *dirCache,
	itemBuilder func(titlePath string, subEntries []os.DirEntry) (T, bool),
	groupBuilder func(name string, items []T) G,
	concurrency int,
) []G {
	var mu sync.Mutex
	groupMap := make(map[string][]T)

	// Count all group directories first for progress bar
	var totalGroups int
	for _, root := range roots {
		if entries := cache.Read(root); entries != nil {
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

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency) // 🧵 limit goroutines

	for _, root := range roots {
		entries := cache.Read(root)
		if entries == nil {
			continue
		}

		for _, group := range entries {
			if !group.IsDir() {
				continue
			}

			wg.Add(1)

			go func(groupName, groupPath string) {
				defer wg.Done()

				sem <- struct{}{}        // 🛑 acquire
				defer func() { <-sem }() // ✅ release

				groupEntries := cache.Read(groupPath)
				if groupEntries == nil {
					_ = bar.Add(1)
					return
				}

				var items []T
				for _, entry := range groupEntries {
					if !entry.IsDir() {
						continue
					}

					titlePath := filepath.Join(groupPath, entry.Name())
					subEntries := cache.Read(titlePath)
					if subEntries == nil {
						continue
					}

					if item, ok := itemBuilder(titlePath, subEntries); ok {
						items = append(items, item)
					}
				}

				if len(items) > 0 {
					mu.Lock()
					groupMap[groupName] = append(groupMap[groupName], items...)
					mu.Unlock()
				}

				_ = bar.Add(1)

			}(group.Name(), filepath.Join(root, group.Name()))
		}
	}

	wg.Wait()
	bar.Finish()

	var groupNames []string
	for name := range groupMap {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	var groups []G
	for _, name := range sortedKeys(groupMap) {
		groups = append(groups, groupBuilder(name, groupMap[name]))
	}
	return groups
}

func scanMedia(
	roots []string,
	cache *dirCache,
	itemBuilder func(titlePath string, subEntries []os.DirEntry) (model.MediaEntry, bool),
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
	for _, root := range roots {
		if entries := cache.Read(root); entries != nil {
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

	for _, root := range roots {
		groupDirs := cache.Read(root)
		if groupDirs == nil {
			continue
		}

		for _, group := range groupDirs {
			if !group.IsDir() {
				continue
			}

			groupPath := filepath.Join(root, group.Name())

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

					if item, ok := itemBuilder(itemPath, subEntries); ok {
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

	// Optional: sort results alphabetically
	sort.Slice(results, func(i, j int) bool {
		return results[i].GetName() < results[j].GetName()
	})

	return results
}
