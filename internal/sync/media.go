package sync

import (
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

func SyncMedia(
	contentType, inputPath string,
	iconMap map[string]*model.SyncedIconEntry,
) model.SyncedOutput {
	raw := model.RawOutput{}
	util.LoadJSON(inputPath, &raw)

	return model.SyncedOutput{
		Type:        "synced",
		GeneratedAt: time.Now(),
		Data:        buildGenres(raw.Data, contentType, iconMap),
	}
}

func buildGenres(
	rawGenres []model.RawGenre,
	contentType string,
	iconMap map[string]*model.SyncedIconEntry,
) []model.SyncedGenre {
	var genres []model.SyncedGenre

	for _, block := range rawGenres {
		genre := model.SyncedGenre{
			Name:  block.Name,
			Items: buildItems(block.Items, contentType, iconMap),
		}
		genres = append(genres, genre)
	}

	return genres
}

func buildItems(
	rawItems []model.RawEntry,
	contentType string,
	iconMap map[string]*model.SyncedIconEntry,
) []model.SyncedItem {
	var items []model.SyncedItem

	for _, item := range rawItems {
		// Skip only if using an alt icon (optional)
		if item.Icon != nil && isAltVariant(item.Icon.ID) {
			continue
		}

		icon := iconMetaFromLocal(item)
		source := iconMetaFromSynced(resolveIconLink(item.Icon, item.Name, item.Path, contentType, iconMap))

		var children []model.SyncedItem
		if item.Items != nil {
			if len(item.Items.Entries) > 0 {
				children = buildItems(item.Items.Entries, contentType, iconMap)
			} else if len(item.Items.Seasons) > 0 {
				for _, season := range item.Items.Seasons {
					children = append(children, model.SyncedItem{
						Type: "season",
						Name: season,
						Path: item.Path + "/" + season,
					})
				}
			}
		}

		items = append(items, model.SyncedItem{
			Type:   item.Type,
			Name:   item.Name,
			Path:   item.Path,
			Status: item.Status,
			Icon:   icon,
			Source: source,
			Items:  children,
		})
	}

	return items
}

func resolveIconLink(
	icon *model.IconMeta,
	name, path, contentType string,
	iconMap map[string]*model.SyncedIconEntry,
) *model.SyncedIconEntry {
	if icon == nil {
		return nil
	}
	baseID := normalizeID(icon.ID)
	linked := iconMap[baseID]
	if linked != nil {
		linked.UsedBy = &model.UsedBy{
			Name:        name,
			Path:        path,
			ContentType: contentType,
		}
	}
	return linked
}

func convertChildren(input any) []model.SyncedItem {
	rawChildren, ok := input.([]model.RawEntry)
	if !ok {
		return nil
	}
	var out []model.SyncedItem
	for _, c := range rawChildren {
		out = append(out, model.SyncedItem{
			Type:   c.Type,
			Name:   c.Name,
			Path:   c.Path,
			Status: c.Status,
		})
	}
	return out
}

func iconMetaFromLocal(item model.RawEntry) *model.IconMeta {
	if item.Icon == nil {
		return nil
	}

	return &model.IconMeta{
		Name:     item.Icon.Name,
		FullPath: item.Icon.FullPath,
		Size:     item.Icon.Size,
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
