package organize

import (
	"strings"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

type MediaSlugMap map[string]model.MovieEntry

func BuildSlugMap(entries []model.MovieEntry) MediaSlugMap {
	slugMap := make(MediaSlugMap)
	for _, e := range entries {
		if e.Icon == nil || e.Icon.Name == "" {
			continue // skip entries with no icon
		}

		slug := util.Slugify(e.Icon.Name)
		slugMap[slug] = e
	}
	return slugMap
}

func MatchIconToEntry(icon model.IconEntry, slugMap MediaSlugMap) (model.MovieEntry, bool) {
	name := strings.TrimSuffix(icon.Name, ".ico")
	slug := util.Slugify(name)

	entry, ok := slugMap[slug]
	return entry, ok
}
