package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/internal/indexer"
)

func main() {
	fmt.Println("🔍 Building .ico index from personal/ and downloaded/ icons...")

	index, err := indexer.BuildIconIndex()
	if err != nil {
		log.Fatalf("❌ Failed to build icon index: %v", err)
	}

	err = indexer.SaveIconIndex(index)
	if err != nil {
		log.Fatalf("❌ Failed to save icon index: %v", err)
	}

	fmt.Printf("✅ ICO index written to %s (%d entries)\n", indexer.OutputPath, len(index.Groups))
}
