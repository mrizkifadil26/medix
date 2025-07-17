package main

import (
	"log"

	"github.com/mrizkifadil26/medix/iconmap"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	flags := iconmap.ParseFlags()

	var cfg iconmap.Config

	if flags.ConfigPath != "" {
		// Load config from JSON file
		if err := utils.LoadJSON(flags.ConfigPath, &cfg); err != nil {
			log.Fatalf("❌ Failed to load config: %v", err)
		}
		log.Println("📄 Loaded config from file:", flags.ConfigPath)
	} else {
		// Inline CLI config
		cfg = iconmap.Config{
			Sources: []iconmap.SourceConfig{
				{
					Path: flags.Source,
					Type: flags.Type,
					Name: flags.Name,
				},
			},
			OutputPath:  flags.Output,
			ExcludeDirs: []string{"done", "New", "todo", "torrent done"},
		}
		log.Println("⚙️ Using inline CLI config")
	}

	index, err := iconmap.GenerateIndex(cfg)
	if err != nil {
		log.Fatalf("❌ Indexing failed: %v", err)
	}

	if err := utils.WriteJSON(cfg.OutputPath, index); err != nil {
		log.Fatalf("❌ Failed to write JSON: %v", err)
	}

	log.Println("✅ Icon index created:", cfg.OutputPath)
}
