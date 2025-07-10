package sync

import (
	"fmt"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

func SyncIcons(path string) (
	*model.SyncedIconIndex,
	map[string]*model.SyncedIconEntry,
) {
	var raw model.SyncedIconIndex
	util.LoadJSON(path, &raw)

	iconMap := ParseIconIndex(&raw)

	raw.Type = "synced-index"
	raw.GeneratedAt = time.Now()

	return &raw, iconMap
}

// ParseIconIndex splits base icons and alt variants.
func ParseIconIndex(raw *model.SyncedIconIndex) map[string]*model.SyncedIconEntry {
	iconMap := make(map[string]*model.SyncedIconEntry)
	variantBuffer := make(map[string][]model.SyncedIconMeta)

	for _, group := range raw.Data {
		for i := range group.Items {
			icon := &group.Items[i]
			baseID := normalizeID(icon.ID)

			// 🔁 If this is a variant (-alt or -alt-N)
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
					// fmt.Printf("🟡 Variant: %-20s → attached to base '%s'\n", icon.ID, baseID)
				} else {
					variantBuffer[baseID] = append(variantBuffer[baseID], meta)
					// fmt.Printf("🕒 Buffered: %-20s → waiting for base '%s'\n", icon.ID, baseID)
				}
				continue // ✅ Don't store alt icon as top-level
			}

			// ✅ Base icon — store it in the map
			iconMap[baseID] = icon
			// fmt.Printf("🟢 Base Icon: %-20s → added\n", icon.ID)

			// 🔁 Attach any previously buffered variants
			if buffered, ok := variantBuffer[baseID]; ok {
				icon.Variants = append(icon.Variants, buffered...)
				// fmt.Printf("🔁 Buffered Variants for '%s': %d attached\n", baseID, len(buffered))
				delete(variantBuffer, baseID)
			}
		}
	}

	// ❗ Report orphaned variants (base icon never appeared)
	for baseID, buffered := range variantBuffer {
		for _, v := range buffered {
			fmt.Printf("🔴 Orphan Variant: %-20s → base '%s' never appeared\n", v.ID, baseID)
		}
	}

	return iconMap
}

func FilterVariants(index *model.SyncedIconIndex) {
	for gi := range index.Data {
		filtered := make([]model.SyncedIconEntry, 0, len(index.Data[gi].Items))
		for _, icon := range index.Data[gi].Items {
			if normalizeID(icon.ID) == icon.ID {
				filtered = append(filtered, icon)
			}
		}
		index.Data[gi].Items = filtered
	}
}
