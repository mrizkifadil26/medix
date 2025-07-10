package main

import (
	"encoding/json"
	"flag"
	"strings"

	"fmt"
	"log"
	"os"

	"github.com/mrizkifadil26/medix/internal/scan"
	"github.com/mrizkifadil26/medix/util"
)

func main() {
	configPath := flag.String("config", "", "Path to scan_config.json (required)")
	filterType := flag.String("type", "", "Optional content type to filter (e.g. movies or tvshows)")

	flag.Parse()

	// Validate required flag
	if *configPath == "" {
		log.Fatalf("❌ Missing required -config flag")
	}

	// Read JSON config
	configBytes, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	var config scan.ScanConfigFile
	if err := json.Unmarshal(configBytes, &config); err != nil {
		log.Fatalf("❌ Failed to parse config JSON: %v", err)
	}

	if len(config.Configs) == 0 {
		log.Fatalf("❌ No configurations found in %s", *configPath)
	}

	var found bool
	for _, cfg := range config.Configs {
		if *filterType != "" && strings.ToLower(cfg.ContentType) != strings.ToLower(*filterType) {
			continue
		}

		found = true
		fmt.Printf("🚀 Scanning %s...\n", cfg.ContentType)
		result := scan.ScanAll(cfg)

		if len(result.Data) == 0 {
			log.Printf("⚠️ No entries found for %s\n", cfg.ContentType)
			continue
		}

		if err := util.WriteJSON(cfg.OutputPath, result); err != nil {
			log.Fatalf("❌ Failed to write JSON for %s: %v", cfg.ContentType, err)
		}
		fmt.Printf("✅ %s written (%d genres)\n", cfg.OutputPath, len(result.Data))
	}

	if !found {
		log.Printf("⚠️ No matching scan config found for -type=%s\n", *filterType)
	}
	// var rootDir, outputPath string

	// switch mode {
	// case "movies":
	// 	rootDir = "/mnt/e/Media/Movies"
	// 	outputPath = "data/movies.raw.json"
	// case "tvshows":
	// 	rootDir = "/mnt/e/Media/TV Shows"
	// 	outputPath = "data/tv_shows.raw.json"
	// default:
	// 	log.Fatalf("Unknown mode: %s", mode)
	// }

	// result := scan.ScanDirectory(mode, rootDir)
	// if len(result.Data) == 0 {
	// 	log.Printf("⚠️ No entries found in %s\n", rootDir)
	// 	return
	// }

	// if err := util.WriteJSON(outputPath, result); err != nil {
	// 	log.Fatalf("❌ Failed to write JSON: %v", err)
	// }
	// fmt.Printf("✅ %s written (%d genres)\n", outputPath, len(result.Data))
}
