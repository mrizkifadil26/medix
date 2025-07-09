package main

import (
	"fmt"
	"os"

	"github.com/mrizkifadil26/medix/internal/sync"
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

	iconIndex := sync.LoadIconIndex("data/ico.index.json")
	iconMap := sync.MapIconsByID(iconIndex)

	// syncAndWrite("movies", moviesPath, outMoviesSynced, iconMap)
	movies := sync.SyncMedia("movies", moviesPath, iconMap)
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
