package main

import (
	"fmt"
	"log"

	scannerV2 "github.com/mrizkifadil26/medix/scanner-v2"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	cli := scannerV2.ParseCLI()

	// CLI scan
	if cli.Root != "" {
		fmt.Println("📂 Scanning using CLI...")
		results, err := scannerV2.Scan(cli.Root, cli.ToOptions())
		if err != nil {
			log.Fatalf("❌ CLI scan failed: %v", err)
		}

		if cli.OutputPath != "" {
			utils.WriteJSON(cli.OutputPath, results)
		}
	}

	// Config file scan
	if cli.ConfigPath != "" {
		cfg, err := scannerV2.LoadConfig(cli.ConfigPath)
		if err != nil {
			log.Fatalf("❌ Failed to load config file: %v", err)
		}

		// Apply CLI overrides
		// cli.OverrideConfig(cfg)

		fmt.Println("📄 Scanning using config file...")
		results, err := scannerV2.Scan(cfg.Root, cfg.ToOptions())
		if err != nil {
			log.Fatalf("❌ Config scan failed: %v", err)
		}

		if cli.OutputPath != "" {
			utils.WriteJSON(cli.OutputPath, results)
		}
	}
}
