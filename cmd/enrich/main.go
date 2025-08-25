package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/enricher"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args, err := enricher.ParseCLI()
	if err != nil {
		log.Fatalf("Error parsing CLI: %v", err)
	}

	var config enricher.Config
	if args.ConfigPath != nil {
		if err := utils.LoadJSON(*args.ConfigPath, &config); err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	fmt.Println("✨ Enriching entries via TMDb...")
	var data any
	if err := utils.LoadJSON(config.Root, &data); err != nil {
		panic(err)
	}

	enriched, err := enricher.Enrich(data, config)
	if err != nil {
		log.Fatalf("❌ Enrichment failed: %v", err)
	}

	fmt.Println("💾 Writing output to:", config.Output)
	if err := utils.WriteJSON(config.Output, enriched); err != nil {
		log.Fatalf("❌ Failed to save output: %v", err)
	}

	fmt.Println("✅ Done enriching.")
}
