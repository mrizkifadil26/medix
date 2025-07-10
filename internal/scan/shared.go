package scan

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func sortedKeys[M ~map[string]V, V any](m M) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func scanGenericGroup[T any, G any](
	roots []string,
	cache *dirCache,
	itemBuilder func(titlePath string, subEntries []os.DirEntry) (T, bool),
	groupBuilder func(name string, items []T) G,
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
