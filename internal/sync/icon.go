package sync

import (
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

func MapIconsByID(index *model.SyncedIconIndex) map[string]*model.SyncedIconEntry {
	result := make(map[string]*model.SyncedIconEntry)
	for _, group := range index.Data {
		for i := range group.Items {
			id := group.Items[i].ID
			if id != "" {
				result[id] = &group.Items[i]
			}
		}
	}
	return result
}
