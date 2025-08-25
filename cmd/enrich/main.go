package main

import (
	"fmt"
	"log"
	"os"

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

	var data any

	// Use existing output as main data if it exists
	if _, err := os.Stat(config.Output); err == nil {
		fmt.Println("⚡ Loading existing enriched data from:", config.Output)
		if err := utils.LoadJSON(config.Output, &data); err != nil {
			log.Fatalf("❌ Failed to load existing output: %v", err)
		}
	} else if os.IsNotExist(err) {
		fmt.Println("✨ No existing output found. Loading root data for enrichment...")
		if err := utils.LoadJSON(config.Root, &data); err != nil {
			log.Fatalf("❌ Failed to load root data: %v", err)
		}
	} else {
		log.Fatalf("❌ Failed to check output file: %v", err)
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
