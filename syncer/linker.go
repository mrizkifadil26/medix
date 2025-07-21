package syncer

import (
	"strconv"

	"github.com/mrizkifadil26/medix/model"
)

// Link links media entries with icon entries based on slugified names.
func Link(media []model.MediaEntry, icons []model.IconEntry) ([]*SyncedMediaEntry, []*SyncedIconEntry) {
	var syncedMedia []*SyncedMediaEntry
	var syncedIcons []*SyncedIconEntry

	// Create reverse map of media that references icons
	iconIDSizeToMedia := make(map[string]*model.MediaEntry)

	for i := range media {
		m := &media[i]
		synced := &SyncedMediaEntry{MediaEntry: *m}

		if m.Icon != nil {
			key := m.Icon.ID + "|" + strconv.FormatInt(m.Icon.Size, 10)
			for j := range icons {
				icon := &icons[j]
				if icon.ID == m.Icon.ID && icon.Size == m.Icon.Size {
					synced.IconSource = icon
					iconIDSizeToMedia[key] = m
					break
				}
			}
		}

		syncedMedia = append(syncedMedia, synced)
	}

	// Now match each icon to the media that used it
	for i := range icons {
		icon := &icons[i]
		synced := &SyncedIconEntry{IconEntry: *icon}

		if icon.ID != "" && icon.Size > 0 {
			key := icon.ID + "|" + strconv.FormatInt(icon.Size, 10)
			if m, ok := iconIDSizeToMedia[key]; ok {
				synced.UsedBy = m
			}
		}

		syncedIcons = append(syncedIcons, synced)
	}

	return syncedMedia, syncedIcons
}
