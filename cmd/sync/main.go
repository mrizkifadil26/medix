package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrizkifadil26/medix/syncer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	configPath := flag.String("config", "", "Path to the sync configuration file (required)")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "❌ Missing required --config argument.")
		flag.Usage()
		os.Exit(1)
	}

	// Load sync configuration
	var cfg syncer.SyncConfig
	err := utils.LoadJSON(*configPath, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to load config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔄 [%s] Syncing using config: %s\n", cfg.Name, *configPath)

	// Load icon index and build map
	icons, err := syncer.LoadIcons(cfg.IconInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to load icons from %s: %v\n", cfg.IconInput, err)
		os.Exit(1)
	}

	entries, err := syncer.LoadMedia(cfg.MediaInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to load media from %s: %v\n", cfg.MediaInput, err)
		os.Exit(1)
	}

	// Sync media
	syncedMedia, syncedIcons := syncer.Link(entries, icons)

	// Write synced media
	err = utils.WriteJSON(cfg.OutMedia, syncedMedia)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write synced media to %s: %v\n", cfg.OutMedia, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Media synced to %s\n", cfg.OutMedia)

	// Write synced icon index
	err = utils.WriteJSON(cfg.OutIcon, syncedIcons)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write icon index to %s: %v\n", cfg.OutIcon, err)
		os.Exit(1)
	}
	fmt.Printf("✅ Icon index synced to %s\n", cfg.OutIcon)

	fmt.Println("✅ Sync completed successfully.")
}
