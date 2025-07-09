package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

const (
	moviesPath  = "data/movies.raw.json"
	tvshowsPath = "data/tvshows.raw.json"
	iconsPath   = "data/ico.index.json"

	outMoviesSynced = "data/movies.synced.json"
	outTVSynced     = "data/tvshows.synced.json"
	outIconSynced   = "data/ico-index.synced.json"
)

func main() {
	fmt.Println("🔄 Syncing icon index with media entries...")

	iconIndex := loadIconIndex()
	iconMap := mapIconsBySlug(iconIndex)
	// syncAndWrite("movies", moviesPath, outMoviesSynced, iconMap)
	movies := syncMedia("movies", moviesPath, iconMap)
	err := util.WriteJSON(outMoviesSynced, movies)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write synced output: %v\n", err)
		os.Exit(1)
	}
	// syncAndWrite("tvshows", tvshowsPath, outTVSloadIconIndexynced, iconMap)

	err = util.WriteJSON(outIconSynced, iconIndex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write icon index: %v\n", err)
		os.Exit(1)
	}

	printUnusedIcons(iconIndex)

	fmt.Println("✅ Sync complete.")
}

func syncMedia(contentType, inputPath string, iconMap map[string]*model.SyncedIconEntry) model.SyncedOutput {
	raw := model.RawOutput{}
	mustLoadJSON(inputPath, &raw)

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

func convertChildren(input []model.RawEntry) []model.SyncedChildItem {
	var out []model.SyncedChildItem
	for _, c := range input {
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
		ID:       item.Icon.ID, // Use ID from the icon if available
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

func mapIconsBySlug(index *model.SyncedIconIndex) map[string]*model.SyncedIconEntry {
	result := make(map[string]*model.SyncedIconEntry)
	for _, g := range index.Data {
		for i := range g.Items {
			id := g.Items[i].ID
			if id != "" {
				result[id] = &g.Items[i]
			}
		}
	}
	return result
}

func loadIconIndex() *model.SyncedIconIndex {
	var out model.SyncedIconIndex
	mustLoadJSON(iconsPath, &out)
	out.Type = "synced-index"
	out.GeneratedAt = time.Now()

	return &out
}

func mustLoadJSON(path string, v any) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(v); err != nil {
		panic(err)
	}
}

// TODO: Implement a more sophisticated unused icon detection
func printUnusedIcons(index *model.SyncedIconIndex) {
	fmt.Println("\n🧹 Unused icons:")
	count := 0
	for _, group := range index.Data {
		for _, entry := range group.Items {
			if entry.UsedBy == nil {
				fmt.Printf("⚠️  %s (%s)\n", entry.Name, entry.FullPath)
				count++
			}
		}
	}
	if count == 0 {
		fmt.Println("✅ All icons are in use.")
	} else {
		fmt.Printf("🔎 Total unused: %d\n", count)
	}
}
