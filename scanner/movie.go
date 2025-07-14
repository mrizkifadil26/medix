package scanner

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type MovieStrategy struct{}

func (MovieStrategy) Scan(roots []string) (model.MediaOutput, error) {
	start := time.Now()
	cache := &dirCache{}
	concurrency := getConcurrency()

	entries := scanMedia(
		roots,
		cache,
		func(folderPath string, dirEntries []os.DirEntry) (model.MediaEntry, bool) {
			group := filepath.Base(filepath.Dir(folderPath)) // genre

			entry := model.MediaEntry{
				BaseEntry: model.BaseEntry{
					Type:   "movie",
					Name:   filepath.Base(folderPath),
					Path:   folderPath,
					Status: resolveStatus(dirEntries),
					Icon:   findIcon(folderPath, dirEntries),
					Group:  group,
				},
			}
			return entry, true
		},
		concurrency,
	)

	// Build group count (unique genre names)
	groupSet := map[string]struct{}{}
	for _, entry := range entries {
		groupSet[entry.Group] = struct{}{}
	}

	output := model.MediaOutput{
		Version:        "1.0.0",
		GeneratedAt:    time.Now(),
		Source:         "movies",
		TotalItems:     len(entries),
		GroupCount:     len(groupSet),
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          entries,
	}

	return output, nil
}

// func extractChildren(parent string, entries []os.DirEntry, cache *dirCache, contentType string) []model.MovieEntry {
// 	// Sort child directories
// 	sort.Slice(entries, func(i, j int) bool {
// 		return entries[i].Name() < entries[j].Name()
// 	})

// 	var children []model.MovieEntry
// 	for _, e := range entries {
// 		if !e.IsDir() {
// 			continue
// 		}

// 		childPath := filepath.Join(parent, e.Name())
// 		subEntries := cache.Read(childPath)
// 		if subEntries == nil {
// 			continue
// 		}

// 		childType := "single"
// 		for _, sub := range subEntries {
// 			if sub.IsDir() {
// 				childType = "collection"
// 				break
// 			}
// 		}

// 		status := resolveStatus(subEntries)
// 		ico := findIcon(childPath, subEntries)
// 		children = append(children, model.MovieEntry{
// 			BaseEntry: model.BaseEntry{
// 				Type:   childType,
// 				Name:   e.Name(),
// 				Path:   childPath,
// 				Status: status,
// 				Icon:   ico,
// 			},
// 		})
// 	}

// 	return children
// }
