package main

import (
	"encoding/json"
	"flag"
	"runtime"
	"strings"

	"fmt"
	"log"
	"os"

	"github.com/mrizkifadil26/medix/internal/scan"
	"github.com/mrizkifadil26/medix/model"
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

	// Determine concurrency (global)
	concurrency := config.Concurrency
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	scan.SetConcurrency(concurrency)

	var found bool
	for _, cfg := range config.Configs {
		if *filterType != "" && strings.ToLower(cfg.ContentType) != strings.ToLower(*filterType) {
			continue
		}

		found = true
		fmt.Printf("🚀 Scanning %s...\n", cfg.ContentType)
		switch strings.ToLower(cfg.ContentType) {
		case "movies":
			runScan(cfg, scan.MovieStrategy{})
		case "tvshows":
			runScan(cfg, scan.TVShowStrategy{})
		default:
			log.Printf("⚠️ Unsupported content type: %s\n", cfg.ContentType)
		}

		if !found {
			log.Printf("⚠️ No matching scan config found for -type=%s\n", *filterType)
		}
	}

	if !found {
		log.Printf("⚠️ No matching scan config found for -type=%s\n", *filterType)
	}
}

func runScan[T any](cfg scan.ScanConfig, strategy scan.ScanStrategy[T]) {
	result := scan.ScanAll(cfg, strategy)

	// This assumes both model.MovieOutput and model.TVShowOutput have .Data
	var dataLen int
	switch v := any(result).(type) {
	case model.MovieOutput:
		dataLen = len(v.Data)
	case model.TVShowOutput:
		dataLen = len(v.Data)
	default:
		log.Printf("⚠️ Unknown result type for %s\n", cfg.ContentType)
		return
	}

	if dataLen == 0 {
		log.Printf("⚠️ No entries found for %s\n", cfg.ContentType)
		return
	}

	if err := util.WriteJSON(cfg.OutputPath, result); err != nil {
		log.Fatalf("❌ Failed to write JSON for %s: %v", cfg.ContentType, err)
	}
	fmt.Printf("✅ %s written (%d groups)\n", cfg.OutputPath, dataLen)
}
