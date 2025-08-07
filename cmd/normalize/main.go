package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrizkifadil26/medix/normalizer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args, err := normalizer.ParseCLI()
	if err != nil {
		log.Fatalf("❌ CLI error: %v", err)
	}

	fmt.Println("Input path:", args.Input)

	var input any
	err = utils.LoadJSON(args.Input, &input)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	// fmt.Println(input)
	// Load config from JSON file if provided
	// if args.ConfigPath != "" {
	// 	fileConfig, err := utils.LoadConfig[normalizer.Config](args.ConfigPath)
	// 	if err != nil {
	// 		log.Fatalf("Failed to load config file: %v", err)
	// 	}

	// 	// Deep merge file config with CLI overrides
	// 	merged, err := utils.MergeDeep(fileConfig, argConfig)
	// 	if err != nil {
	// 		log.Fatalf("Failed to merge config: %v", err)
	// 	}

	// 	args.Config = merged
	// }

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

	// norm := normalizer.New()
	// output, err := norm.Run(input, steps)
	// if err != nil {
	// 	fail("Normalization failed", err)
	// }
	ContinueOnError := true
	registry := normalizer.NewOperators()
	result, err := normalizer.Process(
		input,
		args.Config.Fields,
		registry,
		normalizer.ErrorHandlingOptions{
			ContinueOnError: ContinueOnError,
			CollectErrors:   true,
		},
	)

	// if err != nil {
	// 	// log.Printf("Process failed: %v", err)
	// 	fail("Process failed", err)
	// }

	// if args.OutputPath != "" {
	// 	if err := utils.WriteJSON(args.OutputPath, result); err != nil {
	// 		log.Fatalf("Failed to write output: %v", err)
	// 	}
	// }

	// Always try to write output, even if errors occurred
	if args.OutputPath != "" {
		if err != nil {
			if m, ok := result.(map[string]any); ok {
				if _, exists := m["_errors"]; !exists {
					m["_errors"] = err.Error()
				}
			}
		}

		if writeErr := utils.WriteJSON(args.OutputPath, result); writeErr != nil {
			log.Fatalf("Failed to write output: %v", writeErr)
		}
	}

	if err != nil {
		if ContinueOnError {
			fmt.Println("✅ Process completed with errors. Check output for details.")
		} else {
			fail("❌ Process failed", err)
		}
	}

}

func fail(msg string, err error) {
	fmt.Fprintf(os.Stderr, "❌ %s: %v\n", msg, err)
	os.Exit(1)
}
