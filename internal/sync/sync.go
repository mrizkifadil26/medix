package sync

import (
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

func SyncMedia(contentType, inputPath string, iconMap map[string]*model.SyncedIconEntry) model.SyncedOutput {
	raw := model.RawOutput{}
	util.LoadJSON(inputPath, &raw)

	out := model.SyncedOutput{
		Type:        "synced",
		GeneratedAt: time.Now(),
	}

	for _, block := range raw.Data {
		genre := model.SyncedGenre{
			Name: block.Name,
		}

		for _, item := range block.Items {
			var iconID string
			if item.Icon != nil {
				iconID = item.Icon.ID
			}

			icon := iconMetaFromLocal(item)
			var linked *model.SyncedIconEntry
			if iconID != "" {
				linked = iconMap[iconID]
				if linked != nil {
					linked.UsedBy = &model.UsedBy{
						Name:        item.Name,
						Path:        item.Path,
						ContentType: contentType,
					}
				}
			}

			var children []model.SyncedChildItem
			if rawChildren, ok := item.Items.([]model.RawEntry); ok {
				children = convertChildren(rawChildren)
			}

			genre.Items = append(genre.Items, model.SyncedItem{
				Type:   item.Type,
				Name:   item.Name,
				Path:   item.Path,
				Status: item.Status,
				Icon:   icon,
				Source: iconMetaFromSynced(linked),
				Items:  children,
			})
		}

		out.Data = append(out.Data, genre)
	}

	return out
}

func convertChildren(input any) []model.SyncedChildItem {
	rawChildren, ok := input.([]model.RawEntry)
	if !ok {
		return nil
	}
	var out []model.SyncedChildItem
	for _, c := range rawChildren {
		out = append(out, model.SyncedChildItem{
			Type:   c.Type,
			Name:   c.Name,
			Path:   c.Path,
			Status: c.Status,
		})
	}
	return out
}

func iconMetaFromLocal(item model.RawEntry) *model.SyncedIconMeta {
	if item.Icon == nil {
		return nil
	}
	return &model.SyncedIconMeta{
		Name:     item.Icon.Name,
		FullPath: item.Icon.FullPath,
		Size:     item.Icon.Size,
		Type:     "icon",
		ID:       item.Icon.ID,
	}
}

func iconMetaFromSynced(entry *model.SyncedIconEntry) *model.SyncedIconMeta {
	if entry == nil {
		return nil
	}
	return &model.SyncedIconMeta{
		ID:       entry.ID,
		Name:     entry.Name,
		Size:     entry.Size,
		Source:   entry.Source,
		FullPath: entry.FullPath,
		Type:     entry.Type,
	}
}
