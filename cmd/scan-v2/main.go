package main

import (
	"encoding/json"
	"fmt"
	"log"

	scannerV2 "github.com/mrizkifadil26/medix/scanner-v2"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args := scannerV2.ParseCLI()
	argConfig := args.Config // CLI-level config overrides

	// Load from config file if provided
	if args.ConfigPath != "" {
		fileConfig, err := utils.LoadConfig[scannerV2.Config](args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}

		// Deep merge file config with CLI overrides
		merged, err := utils.MergeDeep(fileConfig, argConfig)
		if err != nil {
			log.Fatalf("Failed to merge config: %v", err)
		}

		// Apply CLI overrides
		// cli.OverrideConfig(cfg)

		fmt.Println("📄 Scanning using config file...")
		PrettyJSON(cfg.ToOptions())
		results, err := scannerV2.Scan(cfg.Root, cfg.ToOptions())
		if err != nil {
			log.Fatalf("❌ Config scan failed: %v", err)
		}

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
