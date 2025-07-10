package scan

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type MovieStrategy struct{}

func (MovieStrategy) Scan(roots []string) model.MovieOutput {
	// result := model.MovieOutput{
	// 	Type:        "raw",
	// 	GeneratedAt: time.Now(),
	// }

	// cache := &dirCache{}
	// groupMap := make(map[string]model.MovieGroup)
	// for _, root := range roots {
	// 	scanMovieRoot(root, groupMap, cache)
	// }
	// // scanMovieRoot(root, groupMap, cache)

	// var groupNames []string
	// for name := range groupMap {
	// 	groupNames = append(groupNames, name)
	// }
	// sort.Strings(groupNames)

	// for _, name := range groupNames {
	// 	result.Data = append(result.Data, groupMap[name])
	// }

	// fmt.Println("Movies scanned.")
	// return result
	cache := &dirCache{}
	groups := scanGenericGroup[model.MovieEntry, model.MovieGroup](
		roots,
		cache,
		func(titlePath string, subEntries []os.DirEntry) (model.MovieEntry, bool) {
			children := extractChildren(titlePath, subEntries, cache, "movies")

			itemType := "single"
			if len(children) > 0 {
				itemType = "collection"
			}

			item := model.MovieEntry{
				BaseEntry: model.BaseEntry{
					Type:   itemType,
					Name:   filepath.Base(titlePath),
					Path:   titlePath,
					Status: resolveStatus(subEntries),
					Icon:   findIcon(titlePath, subEntries),
				},
			}
			if len(children) > 0 {
				item.Items = children
			}

			return item, true
		},
		func(name string, items []model.MovieEntry) model.MovieGroup {
			return model.MovieGroup{Name: name, Items: items}
		},
	)

	return model.MovieOutput{
		Type:        "raw",
		GeneratedAt: time.Now(),
		Data:        groups,
	}
}

// func scanMovieRoot(root string, groupMap map[string]model.MovieGroup, cache *dirCache) {
// 	entries := cache.Read(root)
// 	if entries == nil {
// 		return
// 	}

// 	for _, group := range entries {
// 		if !group.IsDir() {
// 			continue
// 		}

// 		groupName := group.Name()
// 		groupPath := filepath.Join(root, groupName)
// 		groupItems := scanMovieGroup(groupPath, cache)

// 		if len(groupItems) == 0 {
// 			continue
// 		}

// 		groupMap[groupName] = model.MovieGroup{
// 			Name:  groupName,
// 			Items: groupItems,
// 		}
// 	}
// }

// func scanMovieGroup(groupPath string, cache *dirCache) []model.MovieEntry {
// 	entries := cache.Read(groupPath)
// 	if entries == nil {
// 		return nil
// 	}

// 	sort.Slice(entries, func(i, j int) bool {
// 		return entries[i].Name() < entries[j].Name()
// 	})

// 	var items []model.MovieEntry
// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			continue
// 		}

// 		titlePath := filepath.Join(groupPath, entry.Name())
// 		subEntries := cache.Read(titlePath)
// 		if subEntries == nil {
// 			continue
// 		}

// 		children := extractChildren(titlePath, subEntries, cache, "movies")
// 		itemType := "single"
// 		if len(children) > 0 {
// 			itemType = "collection"
// 		}

// 		item := model.MovieEntry{
// 			BaseEntry: model.BaseEntry{
// 				Type:   itemType,
// 				Name:   entry.Name(),
// 				Path:   titlePath,
// 				Status: resolveStatus(subEntries),
// 				Icon:   findIcon(titlePath, subEntries),
// 			},
// 		}

// 		if len(children) > 0 {
// 			item.Type = "collection"
// 			item.Items = children
// 		}

// 		items = append(items, item)
// 	}

// 	return items
// }

// func scanMovieTitle(titlePath string, subEntries []os.DirEntry) (model.MovieEntry, bool) {
// 	children := extractChildren(titlePath, subEntries, cache)
// 	entryType := "single"
// 	if len(children) > 0 {
// 		entryType = "collection"
// 	}

// 	return model.MovieEntry{
// 		BaseEntry: model.BaseEntry{
// 			Type:   entryType,
// 			Name:   filepath.Base(titlePath),
// 			Path:   titlePath,
// 			Status: resolveStatus(subEntries),
// 			Icon:   findIcon(titlePath, subEntries),
// 		},
// 		Items: children,
// 	}, true
// }

func extractChildren(parent string, entries []os.DirEntry, cache *dirCache, contentType string) []model.MovieEntry {
	// Sort child directories
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var children []model.MovieEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		childPath := filepath.Join(parent, e.Name())
		subEntries := cache.Read(childPath)
		if subEntries == nil {
			continue
		}

		childType := "single"
		for _, sub := range subEntries {
			if sub.IsDir() {
				childType = "collection"
				break
			}
		}

		status := resolveStatus(subEntries)
		ico := findIcon(childPath, subEntries)
		children = append(children, model.MovieEntry{
			BaseEntry: model.BaseEntry{
				Type:   childType,
				Name:   e.Name(),
				Path:   childPath,
				Status: status,
				Icon:   ico,
			},
		})
	}

	return children
}
