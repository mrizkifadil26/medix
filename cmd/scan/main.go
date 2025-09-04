package main

import (
	"flag"
	"log"

	"github.com/mrizkifadil26/medix/internal/db"
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

	// Open DB
	database := db.Open("db/sqlite/media.db")
	defer database.Close()

	// results, err := scanner.Scan(
	// 	config.Root,
	// 	*config.Options,
	// 	*config.Output,
	// 	config.Tags,
	// )
	err = scanner.ScanDirectory(database, config.Root, "movie")
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	// Output results
	// outputPath := config.Output.OutputPath
	// if outputPath != "" {
	// 	if err := utils.WriteJSON(outputPath, results); err != nil {
	// 		log.Fatalf("Failed to write output: %v", err)
	// 	}
	// }
}
