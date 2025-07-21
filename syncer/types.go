package syncer

import "github.com/mrizkifadil26/medix/model"

type SyncedMediaEntry struct {
	model.MediaEntry
	IconSource *model.IconEntry `json:"iconSource,omitempty"` // full icon object (not just path)
}

type SyncedIconEntry struct {
	model.IconEntry
	UsedBy *model.MediaEntry `json:"usedBy,omitempty"` // full media object
}
