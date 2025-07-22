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
	filterType := flag.String("type", "", "Filter by type (e.g. media or icon)")
	filterName := flag.String("name", "", "Filter by name (e.g. movies.todo or tv)")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("❌ Missing required -config flag")
	}

	var config scanner.ScanFileConfig
	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("❌ Failed to parse config JSON: %v", err)
	}

	if len(config.Scan) == 0 {
		log.Fatalf("❌ No scan entries in config file")
	}

	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	scanner.SetConcurrency(concurrency)

	var found bool
	for i := range config.Scan {
		cfg := &config.Scan[i]

		if *filterType != "" && !strings.EqualFold(cfg.Type, *filterType) {
			continue
		}

		if *filterName != "" && !strings.EqualFold(cfg.Name, *filterName) {
			continue
		}

		found = true

		var strategy scanner.ScanStrategy

		switch strings.ToLower(cfg.Type) {
		case "movies":
			strategy = scanner.MovieStrategy{}
		case "tv":
			strategy = scanner.TVStrategy{}
		case "icon":
		default:
			log.Printf("⚠️ Skipping unsupported content type: %s\n", cfg.Type)
			continue
		}

		// Prepare scan sources
		sources := make(map[string]string)
		for _, entry := range cfg.Include {
			sources[entry.Label] = entry.Path
		}

		// Apply default options if not set
		if cfg.Options == nil {
			cfg.Options = &scanner.ScanOptions{}
		}

		if cfg.Options.Mode == "" {
			switch strings.ToLower(cfg.Phase) {
			case "raw":
				cfg.Options.Mode = scanner.ScanFiles
			default:
				cfg.Options.Mode = scanner.ScanDirs
			}
		}

		if cfg.Options.Depth == 0 {
			switch strings.ToLower(cfg.Phase) {
			case "raw":
				cfg.Options.Depth = 4
			default:
				cfg.Options.Depth = 2
			}
		}

		if cfg.Options.Concurrency <= 0 {
			cfg.Options.Concurrency = concurrency
		}

		log.Printf("🔍 Scanning...\n")
		log.Printf("⚙️ Options:\n"+
			"   • Name        : %s\n"+
			"   • Type        : %s\n"+
			"   • Phase       : %s\n"+
			"   • Mode        : %s\n"+
			"   • Depth       : %d\n"+
			"   • Concurrency : %d\n",
			cfg.Name,
			cfg.Type,
			cfg.Phase,
			cfg.Options.Mode,
			cfg.Options.Depth,
			cfg.Options.Concurrency)

		output, err := strategy.Scan(sources, *cfg.Options)
		if err != nil {
			log.Printf("❌ Failed to scan %s: %v\n", cfg.Name, err)
			continue
		}

		if len(output.Items) == 0 {
			log.Printf("⚠️ No items found for %s\n", cfg.Name)
			continue
		}

		if err := utils.WriteJSON(cfg.Output, output); err != nil {
			log.Printf("❌ Failed to write output for %s: %v\n", cfg.Name, err)
			continue
		}

		log.Printf("✅ %s written: %d items in %d groups (%d ms)\n",
			cfg.Output, output.TotalItems, output.GroupCount, output.ScanDurationMs)
	}

	if !found {
		log.Printf("⚠️ No matching scan config found for -type=%s\n", *filterType)
	}
}
