package scanner

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type MovieStrategy struct{}

func (MovieStrategy) Scan(sources map[string]string) (model.MediaOutput, error) {
	start := time.Now()
	cache := &dirCache{}
	concurrency := getConcurrency()

	entries := scanMedia(
		sources,
		cache,
		func(folderPath, label string, dirEntries []os.DirEntry) (model.MediaEntry, bool) {
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
		Source:         "movies",
		TotalItems:     len(entries),
		GroupCount:     len(groupSet),
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          entries,
	}

	return output, nil
}
