package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/enricher"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	cfg, err := enricher.ParseCLI()
	if err != nil {
		log.Fatalf("❌ Failed to parse config: %v", err)
	}

	fmt.Println("🔍 Loading input:", cfg.InputFile)
	entries, err := utils.LoadJSONPtr[enricher.Config](cfg.InputFile)
	if err != nil {
		log.Fatalf("❌ Failed to load entries: %v", err)
	}
	// fmt.Printf("📦 Loaded %d entries\n", len(entries))

	// fmt.Println("✨ Enriching entries via TMDb...")
	// enriched, err := Enrich(entries, cfg)
	// if err != nil {
	// 	log.Fatalf("❌ Enrichment failed: %v", err)
	// }

	// fmt.Println("💾 Writing output to:", cfg.OutputFile)
	// if err := utils.WriteJSON(cfg.OutputFile, enriched); err != nil {
	// 	log.Fatalf("❌ Failed to save output: %v", err)
	// }

	fmt.Println("✅ Done enriching.")
}
