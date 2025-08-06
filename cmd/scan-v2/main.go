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
	args := scannerV2.ParseCLI()
	argConfig := args.Config // CLI-level config overrides

	// Load from config file if provided
	if args.ConfigPath != nil {
		fileConfig, err := utils.LoadConfig[scannerV2.Config](*args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
			os.Exit(1)
		}

		// Deep merge file config with CLI overrides
		merged, err := utils.MergeDeep(fileConfig, argConfig)
		if err != nil {
			log.Fatalf("Failed to merge config: %v", err)
		}

		args.Config = merged
	}

	// Validate required field
	if argConfig.Root == nil || *argConfig.Root == "" {
		flag.Usage()
		log.Fatal("Error: --root is required (or must be in config file)")
	}

	args.Config.ApplyDefaults() // Apply defaults to ensure all options are set

	// root := args.Config.Root
	// opts := args.Config.Options

	// output, err := scannerV2.Scan(root, opts)
	// if err != nil {
	// log.Fatal(err)
	// }

	if args.OutputPath != "" {
		if err := utils.WriteJSON(args.OutputPath, output); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
	}
}

func PrettyJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
