package scanner

import (
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type MovieStrategy struct{}

func (MovieStrategy) Scan(sources map[string]string) (model.MediaOutput, error) {
	start := time.Now()
	cache := &dirCache{}
	concurrency := getConcurrency()

	entries := Scan(
		sources,
		cache,
		ScanOptions{
			Mode:         ScanDirs,
			Depth:        2,
			Exts:         []string{},
			Concurrency:  concurrency,
			ShowProgress: true,
		},
		func(item ScannedItem) (model.MediaEntry, bool) {
			folderPath := item.ItemPath
			dirEntries := item.SubEntries
			label := item.GroupLabel

			group := filepath.Base(filepath.Dir(folderPath)) // genre

			entry := model.MediaEntry{
				BaseEntry: model.BaseEntry{
					Type:   "movie",
					Name:   filepath.Base(folderPath),
					Path:   folderPath,
					Status: resolveStatus(dirEntries),
					Icon:   resolveIcon(folderPath, dirEntries),
					Group:  group,
				},
				Source: label,
			}

			return entry, true
		},
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
		Source:         "movies",
		TotalItems:     len(entries),
		GroupCount:     len(groupSet),
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          entries,
	}

	return output, nil
}
