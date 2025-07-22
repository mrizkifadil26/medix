package scanner

import (
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type MovieStrategy struct{}

func (MovieStrategy) Scan(
	sources map[string]string,
	opts ScanOptions,
) (model.MediaOutput, error) {
	start := time.Now()
	cache := &dirCache{}
	concurrency := getConcurrency()

	// Apply defaults if not set
	if opts.Mode == "" {
		opts.Mode = ScanDirs
	}
	if opts.Depth <= 0 {
		opts.Depth = 2
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = concurrency
	}
	opts.ShowProgress = true // always enable for now

	entries := Scan(
		sources,
		cache,
		opts,
		func(item ScannedItem) (model.MediaEntry, bool) {
			folderPath := item.ItemPath
			dirEntries := item.SubEntries
			source := item.Source

			group := item.GroupLabel

			entry := model.MediaEntry{
				BaseEntry: model.BaseEntry{
					Type:   "movie",
					Name:   filepath.Base(folderPath),
					Path:   folderPath,
					Status: resolveStatus(dirEntries),
					Icon:   resolveIcon(folderPath, dirEntries),
					Group:  group,
				},
				Source: source,
			}

			return entry, true
		},
	)

	// Build group count (unique genre names)
	groupSet := map[string]struct{}{}
	for _, entry := range entries {
		if len(entry.Group) == 0 {
			groupSet["<ungrouped>"] = struct{}{}
			continue
		}

		for i := 1; i <= len(entry.Group); i++ {
			groupKey := filepath.Join(entry.Group[:i]...)
			groupSet[groupKey] = struct{}{}
		}
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
