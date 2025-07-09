package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/internal/index"
	"github.com/mrizkifadil26/medix/util"
)

const configPath = "config/ico-indexer.json"

func main() {
	var cfg index.IconIndexerUnifiedConfig
	err := util.LoadJSON(configPath, &cfg)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	for category, target := range cfg.Outputs {
		log.Printf("📦 Indexing category: %s", category)

		indexed, err := index.BuildIconIndex(index.IconIndexerConfig{
			Sources:     target.Sources,
			OutputPath:  target.OutputPath,
			ExcludeDirs: cfg.ExcludeDirs,
		})
		if err != nil {
			log.Printf("⚠️  Failed to index %s: %v", category, err)
			continue
		}

		err = util.WriteJSON(target.OutputPath, indexed)
		if err != nil {
			log.Printf("❌ Failed to write %s index: %v", category, err)
			continue
		}

		fmt.Printf("✅ %s index written to %s\n", category, target.OutputPath)
	}
}
