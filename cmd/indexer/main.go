package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/internal/indexer"
	"github.com/mrizkifadil26/medix/util"
)

func main() {
	cfg := indexer.IconIndexerConfig{
		Sources: []indexer.IconSource{
			{
				Path:   "/mnt/c/Users/Rizki/OneDrive/Pictures/Icons/Personal Icon Pack/Movies/ICO",
				Source: "personal",
			},
			{
				Path:   "/mnt/c/Users/Rizki/OneDrive/Pictures/Icons/Downloaded Icon Pack/Movie Icon Pack/downloaded",
				Source: "downloaded",
			},
		},
		OutputPath:  "data/ico.index.json",
		ExcludeDirs: []string{"Collection"},
	}

	index, err := indexer.BuildIconIndex(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to build icon index: %v", err)
	}

	err = util.WriteJSON(cfg.OutputPath, index)
	if err != nil {
		log.Fatalf("❌ Failed to save icon index: %v", err)
	}

	fmt.Printf("✅ ICO index written to %s (%d entries)\n", cfg.OutputPath, len(index.Data))
}
