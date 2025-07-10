package main

import (
	"fmt"
	"os"

	"github.com/mrizkifadil26/medix/internal/sync"
	"github.com/mrizkifadil26/medix/util"
)

const (
	moviesPath  = "data/movies.raw.json"
	tvshowsPath = "data/tvshows.raw.json"
	iconsPath   = "data/movies.ico.index.json"

	outMoviesSynced = "data/movies.synced.json"
	outTVSynced     = "data/tvshows.synced.json"
	outIconSynced   = "data/ico-index.synced.json"
)

func main() {
	fmt.Println("🔄 Syncing icon index with media entries...")

	iconIndex := sync.LoadIconIndex(iconsPath)
	iconMap := sync.MapIconsByID(iconIndex)

	// syncAndWrite("movies", moviesPath, outMoviesSynced, iconMap)
	movies := sync.SyncMedia("movies", moviesPath, iconMap)
	err := util.WriteJSON(outMoviesSynced, movies)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write synced output: %v\n", err)
		os.Exit(1)
	}
	// syncAndWrite("tvshows", tvshowsPath, outTVSloadIconIndexynced, iconMap)

	iconIndex.Data = sync.FlattenIconMap(iconMap)
	err = util.WriteJSON(outIconSynced, iconIndex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write icon index: %v\n", err)
		os.Exit(1)
	}

	if err := sync.GenerateUnusedIconsReport(iconIndex); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write unused icons report: %v\n", err)
	}
	fmt.Println("✅ Sync complete.")
}
