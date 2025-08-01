package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrizkifadil26/medix/normalizer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args := normalizer.ParseCLI()
	argConfig := args.Config // CLI-level config overrides

	// Load config from JSON file if provided
	if args.ConfigPath != "" {
		fileConfig, err := utils.LoadConfig[normalizer.Config](args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}

		// Deep merge file config with CLI overrides
		merged, err := utils.MergeDeep(fileConfig, argConfig)
		if err != nil {
			log.Fatalf("Failed to merge config: %v", err)
		}

		args.Config = merged
	}

	/*
		 else {
			// Attempt to parse JSON input (array or primitive)
			if strings.HasSuffix(inputFlag, ".json") {
				var scan model.MediaOutput
				if err := utils.LoadJSON(inputFlag, &scan); err != nil {
					fail("Failed to load input JSON", err)
				}
				// Extract names into string array
				var names []any
				for _, item := range scan.Items {
					names = append(names, item.Name)
				}
				input = names
			} else {
				var tryParse any
				if err := json.Unmarshal([]byte(inputFlag), &tryParse); err == nil {
					input = tryParse
				} else {
					input = inputFlag // fallback: plain string
				}
			}
			steps = splitAndTrim(stepsFlag)
		}
	*/

	// TODO: Handle validation of CLI
	input := args.Input
	steps := args.Config.Steps

	norm := normalizer.New()
	output, err := norm.Run(input, steps)
	if err != nil {
		fail("Normalization failed", err)
	}

	if args.OutputPath != "" {
		if err := utils.WriteJSON(args.OutputPath, output); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
	}
}

func fail(msg string, err error) {
	fmt.Fprintf(os.Stderr, "❌ %s: %v\n", msg, err)
	os.Exit(1)
}
