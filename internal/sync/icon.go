package sync

import (
	"fmt"
	"regexp"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

func LoadIconIndex(path string) *model.SyncedIconIndex {
	var index model.SyncedIconIndex
	util.LoadJSON(path, &index)

	index.Type = "synced-index"
	index.GeneratedAt = time.Now()
	return &index
}

var altIDRegex = regexp.MustCompile(`-alt(?:-\d+)?$`)

func normalizeID(id string) string {
	return altIDRegex.ReplaceAllString(id, "")
}

func MapIconsByID(index *model.SyncedIconIndex) map[string]*model.SyncedIconEntry {
	iconMap := make(map[string]*model.SyncedIconEntry)
	variantBuffer := make(map[string][]model.SyncedIconMeta)

	for _, group := range index.Data {
		for i := range group.Items {
			icon := &group.Items[i]
			baseID := normalizeID(icon.ID)

			// Skip adding alt variants as top-level
			if baseID != icon.ID {
				meta := model.SyncedIconMeta{
					ID:       icon.ID,
					Name:     icon.Name,
					FullPath: icon.FullPath,
					Size:     icon.Size,
					Source:   icon.Source,
					Type:     icon.Type,
				}

				if baseIcon, ok := iconMap[baseID]; ok {
					baseIcon.Variants = append(baseIcon.Variants, meta)
					// fmt.Printf("🟡 Variant:      %-20s → attached to '%s'\n", icon.ID, baseID)
				} else {
					variantBuffer[baseID] = append(variantBuffer[baseID], meta)
					// fmt.Printf("🕒 Buffered:      %-20s → waiting for base '%s'\n", icon.ID, baseID)
				}

				continue // ← skip storing this variant in map
			}

			// Main/base icon
			iconMap[baseID] = icon
			// fmt.Printf("🟢 Base Icon:    %-20s → added\n", icon.ID)

			// Merge in any buffered variants
			if buffered, ok := variantBuffer[baseID]; ok {
				icon.Variants = append(icon.Variants, buffered...)
				// fmt.Printf("🔁 Buffered Variants for '%s': %d attached\n", baseID, len(buffered))
				delete(variantBuffer, baseID)
			}
		}
	}

	// Log unmerged variants
	for baseID, buffered := range variantBuffer {
		for _, v := range buffered {
			fmt.Printf("🔴 Orphan Variant: %-20s → base '%s' never appeared\n", v.ID, baseID)
		}
	}

	return iconMap
}

// FlattenIconMap returns a slice of genre groups from the final icon map.
func FlattenIconMap(iconMap map[string]*model.SyncedIconEntry) []model.SyncedIconGroup {
	genreMap := make(map[string]*model.SyncedIconGroup)

	for _, icon := range iconMap {
		// Fallback if missing
		genre := "Uncategorized"
		if icon.Source != "" {
			genre = icon.Source // Optional: You may want to categorize by genre instead
		}

		if _, ok := genreMap[genre]; !ok {
			genreMap[genre] = &model.SyncedIconGroup{Name: genre}
		}
		genreMap[genre].Items = append(genreMap[genre].Items, *icon)
	}

	var flat []model.SyncedIconGroup
	for _, g := range genreMap {
		flat = append(flat, *g)
	}
	return flat
}
