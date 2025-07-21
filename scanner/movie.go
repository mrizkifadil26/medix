package scanner

import (
	"os"
	"path/filepath"
	"strings"
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

			matchedRoot := ""
			for _, root := range roots {
				if rel, err := filepath.Rel(root, folderPath); err == nil && !strings.HasPrefix(rel, "..") {
					matchedRoot = root
					break
				}
			}
			entry := model.MediaEntry{
				BaseEntry: model.BaseEntry{
					Type:   "movie",
					Name:   filepath.Base(folderPath),
					Path:   folderPath,
					Status: resolveStatus(dirEntries),
					Icon:   resolveIcon(folderPath, dirEntries),
					Group:  group,
				},
				Source: matchedRoot,
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
