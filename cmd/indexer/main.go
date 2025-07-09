package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/internal/indexer"
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

	err = indexer.SaveIconIndex(index)
	if err != nil {
		log.Fatalf("❌ Failed to save icon index: %v", err)
	}

	fmt.Printf("✅ ICO index written to %s (%d entries)\n", indexer.OutputPath, len(index.Groups))
}
