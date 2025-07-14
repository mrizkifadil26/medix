package main

import (
	"encoding/json"
	"flag"
	"runtime"
	"strings"

	"log"
	"os"

	"github.com/mrizkifadil26/medix/scanner"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	configPath := flag.String("config", "", "Path to scan_config.json (required)")
	filterType := flag.String("type", "", "Filter by content type (e.g. movies or tvshows)")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("❌ Missing required -config flag")
	}

	var config scanner.ScanConfigFile
	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("❌ Failed to parse config JSON: %v", err)
	}
	if len(config.Configs) == 0 {
		log.Fatalf("❌ No scan entries in config file")
	}

	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	scanner.SetConcurrency(concurrency)

	var found bool
	for _, cfg := range config.Configs {
		if *filterType != "" && !strings.EqualFold(cfg.ContentType, *filterType) {
			continue
		}
		found = true

		var (
			strategy scanner.ScanStrategy
			name     = strings.ToLower(cfg.ContentType)
		)

		switch name {
		case "movies":
			strategy = scanner.MovieStrategy{}
		case "tvshows":
			strategy = scanner.TVShowStrategy{}
		default:
			log.Printf("⚠️ Skipping unsupported content type: %s\n", cfg.ContentType)
			continue
		}

		log.Printf("🔍 Scanning %s...\n", name)
		output, err := strategy.Scan(cfg.Sources)
		if err != nil {
			log.Printf("❌ Failed to scan %s: %v\n", name, err)
			continue
		}
		if len(output.Items) == 0 {
			log.Printf("⚠️ No items found for %s\n", name)
			continue
		}

		if err := utils.WriteJSON(cfg.OutputPath, output); err != nil {
			log.Printf("❌ Failed to write output for %s: %v\n", name, err)
			continue
		}

		log.Printf("✅ %s written: %d items in %d groups (%d ms)\n",
			cfg.OutputPath, output.TotalItems, output.GroupCount, output.ScanDurationMs)
	}

	if !found {
		log.Printf("⚠️ No matching scan config found for -type=%s\n", *filterType)
	}
}
