package scanner

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type TVStrategy struct{}

func (TVStrategy) Scan(roots []string) (model.MediaOutput, error) {
	start := time.Now()
	cache := &dirCache{}
	concurrency := getConcurrency()

	entries := scanMedia(
		roots,
		cache,
		func(folderPath string, dirEntries []os.DirEntry) (model.MediaEntry, bool) {
			group := filepath.Base(filepath.Dir(folderPath)) // genre

			showEntry := model.MediaEntry{
				BaseEntry: model.BaseEntry{
					Type:   "show",
					Name:   filepath.Base(folderPath),
					Path:   folderPath,
					Status: resolveStatus(dirEntries),
					Icon:   resolveIcon(folderPath, dirEntries),
					Group:  group,
				},
			}

			// Add seasons as items (not recursive)
			for _, entry := range dirEntries {
				if entry.IsDir() {
					seasonPath := filepath.Join(folderPath, entry.Name())
					subEntries := cache.Read(seasonPath)
					season := model.MediaEntry{
						BaseEntry: model.BaseEntry{
							Type:   "season",
							Name:   entry.Name(),
							Path:   seasonPath,
							Status: resolveStatus(subEntries),
							Parent: showEntry.Name,
							Group:  group,
						},
					}
					showEntry.Items = append(showEntry.Items, season)
				}
			}

			return showEntry, true
		},
		concurrency,
	)

	// Build group count (unique genre names)
	groupSet := map[string]struct{}{}
	for _, entry := range entries {
		groupSet[entry.Group] = struct{}{}
	}

	output := model.MediaOutput{
		Type:           "raw",
		Version:        "1.0.0",
		GeneratedAt:    time.Now(),
		Source:         "tv",
		TotalItems:     len(entries),
		GroupCount:     len(groupSet),
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          entries,
	}

	return output, nil
}
