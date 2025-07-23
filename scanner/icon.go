package scanner

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type IconStrategy struct{}

func (IconStrategy) Scan(
	sources map[string]string,
	opts ScanOptions,
) (any, error) {
	start := time.Now()
	cache := &dirCache{}
	concurrency := getConcurrency()

	// Apply defaults if not set
	if opts.Mode == "" {
		opts.Mode = ScanFiles
	}
	if opts.Depth <= 0 {
		opts.Depth = 2
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = concurrency
	}
	opts.ShowProgress = true // always enable for now

	fmt.Println("mode: ", opts.Mode)
	entries := Scan(
		sources,
		cache,
		opts,
		func(item ScanEntry) (model.IconEntry, bool) {
			folderPath := item.ItemPath
			dirEntries := item.SubEntries
			source := item.Source

			group := item.GroupLabel
			var size int64
			if item.ItemSize != nil {
				size = *item.ItemSize
			}

			entry := model.IconEntry{
				BaseEntry: model.BaseEntry{
					Name:        filepath.Base(folderPath),
					Path:        folderPath,
					Type:        "icon",
					ContentType: "movie",
					Source:      source,
					Group:       group,
				},
				Slug: resolveStatus(dirEntries),
				Size: size,
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

	sourceEntries := make([]string, 0, len(sources))
	for k := range sources {
		sourceEntries = append(sourceEntries, k)
	}

	output := model.IconIndex{
		Type:           "raw",
		Version:        "1.0.0",
		GeneratedAt:    time.Now(),
		Sources:        sourceEntries,
		TotalItems:     len(entries),
		GroupCount:     len(groupSet),
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          entries,
	}

	return output, nil
}
