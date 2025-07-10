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
