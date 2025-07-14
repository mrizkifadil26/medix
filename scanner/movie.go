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
	cache := &dirCache{}
	concurrency := getConcurrency()

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
		concurrency,
	)

	return model.MovieOutput{
		Type:        "raw",
		GeneratedAt: time.Now(),
		Data:        groups,
	}
}

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
