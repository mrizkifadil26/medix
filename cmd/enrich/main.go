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

	var data = utils.NewOrderedMap[string, any]()
	loadPath := config.Root

	fmt.Println("⚡ Loading root data for enrichment from:", loadPath)
	if err := utils.LoadJSON(loadPath, data); err != nil {
		log.Fatalf("❌ Failed to load root data from %s: %v", loadPath, err)
	}

	enriched, err := enricher.Enrich(data, &config)
	if err != nil {
		log.Fatalf("❌ Enrichment failed: %v", err)
	}

	fmt.Println("💾 Writing output to:", config.Output)
	if err := utils.WriteJSON(config.Output, enriched); err != nil {
		log.Fatalf("❌ Failed to save output: %v", err)
	}

	fmt.Println("✅ Done enriching.")
}
