package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	scannerV2 "github.com/mrizkifadil26/medix/scanner-v2"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args, err := scannerV2.ParseCLI()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Start config from CLI
	config := args.Config
	// If config file exists, load and merge
	if args.ConfigPath != nil {
		fileConfig, err := utils.LoadConfig[scannerV2.Config](*args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
			os.Exit(1)
		}

		// Deep merge file config with CLI overrides
		merged, err := utils.Merge(
			fileConfig,
			args.Config,
			utils.MergeOptions{
				Overwrite: true,
				Recursive: true,
			},
		)

		if err != nil {
			log.Fatalf("Failed to merge CLI config: %v", err)
		}

		config = merged
	}

	// Validate required field
	if config.Root == nil || *config.Root == "" {
		flag.Usage()
		log.Fatal("Error: --root is required (or must be in config file)")
	}

	// Fill missing defaults
	finalConfig := config.ApplyDefaults()
	finalConfig.PrettyPrint()

	// output, err := scannerV2.Scan(*config.Root, config.Options)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Output results
	// if args.OutputPath != nil && *args.OutputPath != "" {
	// 	if err := utils.WriteJSON(*args.OutputPath, output); err != nil {
	// 		log.Fatalf("Failed to write output: %v", err)
	// 	}
	// }
}

func PrettyJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
