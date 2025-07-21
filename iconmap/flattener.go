package iconmap

import (
	"strings"

	"github.com/mrizkifadil26/medix/model"
)

func FlattenGroupedIcons(grouped map[string][]model.IconEntry, excludeDirs []string) []model.IconEntry {
	var flat []model.IconEntry

	for groupKey, entries := range grouped {
		groupParts := strings.Split(groupKey, "/")
		for _, entry := range entries {
			entry.Group = groupParts
			flat = append(flat, entry)
		}
	}

	return flat
}
