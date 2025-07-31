package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/normalizer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	var (
		inputFlag  string
		outputFlag string
		stepsFlag  string
		configFlag string
	)

	flag.StringVar(&inputFlag, "input", "", "Input string, JSON array, or path to scan JSON")
	flag.StringVar(&outputFlag, "output", "", "Output file (optional)")
	flag.StringVar(&stepsFlag, "steps", "", "Comma-separated normalization steps")
	flag.StringVar(&configFlag, "config", "", "Path to config JSON (optional)")
	flag.Parse()

	var input any
	var steps []string

	// Load config from JSON file if provided
	if configFlag != "" {
		var job normalizer.NormalizeJob
		if err := utils.LoadJSON(configFlag, &job); err != nil {
			fail("Failed to load config", err)
		}
		input = job.Input
		steps = job.Steps

	} else {
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

	norm := normalizer.New()
	result, err := norm.Run(input, steps)
	if err != nil {
		fail("Normalization failed", err)
	}

	if outputFlag != "" {
		if err := utils.WriteJSON(outputFlag, result); err != nil {
			fail("Failed to write output", err)
		}
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
	}
}

func fail(msg string, err error) {
	fmt.Fprintf(os.Stderr, "❌ %s: %v\n", msg, err)
	os.Exit(1)
}

func splitAndTrim(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
