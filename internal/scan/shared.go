package scan

import (
	"os"
	"path/filepath"
	"sort"
)

func sortedKeys[M ~map[string]V, V any](m M) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func sortDirEntries(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}

func scanGenericGroup[T any, G any](
	roots []string,
	cache *dirCache,
	scanItem func(titlePath string, subEntries []os.DirEntry) (T, bool),
	groupBuilder func(name string, items []T) G,
) []G {
	groupMap := make(map[string][]T)

	for _, root := range roots {
		rootEntries := cache.Read(root)
		if rootEntries == nil {
			continue
		}

		for _, group := range rootEntries {
			if !group.IsDir() {
				continue
			}
			groupName := group.Name()
			groupPath := filepath.Join(root, groupName)

			groupEntries := cache.Read(groupPath)
			if groupEntries == nil {
				continue
			}

			sortDirEntries(groupEntries)
			for _, title := range groupEntries {
				if !title.IsDir() {
					continue
				}

				titlePath := filepath.Join(groupPath, title.Name())
				subEntries := cache.Read(titlePath)
				if subEntries == nil {
					continue
				}

				item, ok := scanItem(titlePath, subEntries)
				if ok {
					groupMap[groupName] = append(groupMap[groupName], item)
				}
			}
		}
	}

	var groups []G
	for _, name := range sortedKeys(groupMap) {
		groups = append(groups, groupBuilder(name, groupMap[name]))
	}
	return groups
}
