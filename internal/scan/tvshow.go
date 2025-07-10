package scan

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type TVShowStrategy struct{}

func (TVShowStrategy) Scan(roots []string) model.TVShowOutput {
	// result := model.TVShowOutput{
	// 	Type:        "raw",
	// 	GeneratedAt: time.Now(),
	// }

	// cache := &dirCache{}
	// groupMap := make(map[string]model.TVShowGroup)
	// for _, root := range roots {
	// 	scanTVRoot(root, groupMap, cache)
	// }

	// var groupNames []string
	// for name := range groupMap {
	// 	groupNames = append(groupNames, name)
	// }
	// sort.Strings(groupNames)

	// for _, name := range groupNames {
	// 	result.Data = append(result.Data, groupMap[name])
	// }

	// fmt.Println("TV Shows scanned.")
	// return result
	cache := &dirCache{}
	groups := scanGenericGroup[model.TVShowEntry, model.TVShowGroup](
		roots,
		cache,
		func(titlePath string, subEntries []os.DirEntry) (model.TVShowEntry, bool) {
			seasons := extractSeasonNames(subEntries)

			return model.TVShowEntry{
				BaseEntry: model.BaseEntry{
					Type:   "single",
					Name:   filepath.Base(titlePath),
					Path:   titlePath,
					Status: resolveStatus(subEntries),
					Icon:   findIcon(titlePath, subEntries),
				},
				Seasons: seasons,
			}, true
		},
		func(name string, items []model.TVShowEntry) model.TVShowGroup {
			return model.TVShowGroup{Name: name, Items: items}
		},
	)

	return model.TVShowOutput{
		Type:        "raw",
		GeneratedAt: time.Now(),
		Data:        groups,
	}
}

// func scanTVRoot(root string, groupMap map[string]model.TVShowGroup, cache *dirCache) {
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
// 		groupItems := scanTVGroup(groupPath, cache)

// 		if len(groupItems) == 0 {
// 			continue
// 		}

// 		groupMap[groupName] = model.TVShowGroup{
// 			Name:  groupName,
// 			Items: groupItems,
// 		}
// 	}
// }

// func scanTVGroup(groupPath string, cache *dirCache) []model.TVShowEntry {
// 	entries := cache.Read(groupPath)
// 	if entries == nil {
// 		return nil
// 	}

// 	sort.Slice(entries, func(i, j int) bool {
// 		return entries[i].Name() < entries[j].Name()
// 	})

// 	var items []model.TVShowEntry
// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			continue
// 		}

// 		titlePath := filepath.Join(groupPath, entry.Name())
// 		subEntries := cache.Read(titlePath)
// 		if subEntries == nil {
// 			continue
// 		}

// 		item := model.TVShowEntry{
// 			BaseEntry: model.BaseEntry{
// 				Type:   "single",
// 				Name:   entry.Name(),
// 				Path:   titlePath,
// 				Status: resolveStatus(subEntries),
// 				Icon:   findIcon(titlePath, subEntries),
// 			},
// 		}

// 		seasons := extractSeasonNames(subEntries)
// 		if len(seasons) > 0 {
// 			item.Seasons = seasons
// 		}

// 		items = append(items, item)
// 	}

// 	return items
// }

func extractSeasonNames(entries []os.DirEntry) []string {
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}
