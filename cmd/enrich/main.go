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

	var rawData = utils.NewOrderedMap[string, any]()
	loadPath := config.Root

	// Decide data source
	if !*args.Refresh {
		if _, err := os.Stat(config.Output); err == nil {
			loadPath = config.Output
			fmt.Println("⚡ Loading existing enriched data from:", loadPath)
		} else if os.IsNotExist(err) {
			fmt.Println("✨ No existing output found. Loading root data for enrichment...")
		} else {
			log.Fatalf("❌ Failed to check output file: %v", err)
		}
	} else {
		fmt.Println("🔄 Refresh mode: ignoring existing output, loading root data...")
	}

	if err := utils.LoadJSON(loadPath, rawData); err != nil {
		log.Fatalf("❌ Failed to load data from %s: %v", loadPath, err)
	}

	// enriched, err := enricher.Enrich(rawData, &config)
	if err != nil {
		log.Fatalf("❌ Enrichment failed: %v", err)
	}

	fmt.Println("💾 Writing output to:", config.Output)
	if err := utils.WriteJSON(config.Output, enriched); err != nil {
		log.Fatalf("❌ Failed to save output: %v", err)
	}

	fmt.Println("✅ Done enriching.")
}
