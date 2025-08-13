package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/normalizer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args, err := normalizer.ParseCLI()
	if err != nil {
		log.Fatalf("Error parsing CLI: %v", err)
	}

	var config normalizer.Config
	if args.ConfigPath != nil {
		config, err = utils.LoadConfig[normalizer.Config](*args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	if err := utils.MergeInto(&config, &args.Config, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	}); err != nil {
		log.Fatalf("Failed to merge CLI config: %v", err)
	}

	var input any
	utils.LoadJSON(config.Root, &input)

	ContinueOnError := false
	registry := normalizer.NewOperators()
	result, err := normalizer.Process(
		input,
		config.Fields,
		registry,
		normalizer.ErrorHandlingOptions{
			ContinueOnError: ContinueOnError,
			CollectErrors:   true,
		},
	)

	// Always try to write output, even if errors occurred
	if config.OutputPath != "" {
		if err := utils.WriteJSON(config.OutputPath, result); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
	}

	if err != nil {
		if ContinueOnError {
			fmt.Println("✅ Process completed with errors. Check output for details.")
		} else {
			log.Fatalf("❌ Process failed: %v", err)
		}
	} else {
		fmt.Println("✅ Process completed")
	}
}
