package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/scanner"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	config := scanner.DefaultConfig()

	args, err := scanner.ParseCLI()
	if err != nil {
		log.Fatalf("Error parsing CLI: %v", err)
	}

	if args.ConfigPath != nil {
		fileConfig, err := utils.LoadConfig[scanner.Config](*args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}

		// Deep merge file config with CLI overrides
		if err := utils.MergeInto(
			&config,
			&fileConfig,
			utils.MergeOptions{
				Overwrite: true,
				Recursive: true,
			},
		); err != nil {
			log.Fatalf("Failed to merge file config into defaults: %v", err)
		}
	}

	// 3. Merge CLI config into result (CLI overrides file+defaults)
	if err := utils.MergeInto(&config, &args.Config, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	}); err != nil {
		log.Fatalf("Failed to merge CLI config: %v", err)
	}

	// Validate required field
	if config.Root == "" {
		flag.Usage()
		log.Fatal("Error: --root is required (or must be in config file)")
	}

	// Fill missing defaults
	// if err := config.ApplyDefaults(); err != nil {
	// 	log.Fatalf("Error applying defaults: %v", err)
	// }
	// config.PrettyPrint()

	results, err := scanner.Scan(
		config.Root,
		*config.Options,
		*config.Output,
		config.Tags,
	)
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	// Output results
	outputPath := config.Output.OutputPath
	if outputPath != "" {
		if err := utils.WriteJSON(outputPath, results); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
	}
}

func PrettyJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
