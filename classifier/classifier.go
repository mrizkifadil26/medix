package classifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Run is the main entrypoint to run classification
func Run(cfg Config) error {
	var all []Entry
	var inputFiles []string
	seen := make(map[string]bool)

	// Expand glob patterns in input
	for _, pattern := range cfg.Inputs {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("invalid input pattern %q: %w", pattern, err)
		}
		for _, f := range files {
			if !seen[f] {
				inputFiles = append(inputFiles, f)
				seen[f] = true
			}
		}
	}

	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files matched any pattern: %v", cfg.Inputs)
	}

	// Load entries from each matched input file
	for _, in := range inputFiles {
		data, err := os.ReadFile(in)
		if err != nil {
			return fmt.Errorf("read input %s: %w", in, err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("parse input %s: %w", in, err)
		}

		entriesRaw, ok := result["items"]
		if !ok {
			return fmt.Errorf("missing 'items' in %s", in)
		}

		entries, ok := entriesRaw.([]interface{})
		if !ok {
			return fmt.Errorf("'items' should be an array in %s", in)
		}

		for _, e := range entries {
			entry, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			entry["source"] = in
			all = append(all, entry)
		}
	}

	// Apply output rules
	for _, rule := range cfg.Outputs {
		var out Output
		top := make(map[string]interface{})

		if rule.Set != nil {
			for k, v := range rule.Set {
				top[k] = v
			}
		}

		for _, entry := range all {
			if matchEntry(entry, rule.Match) {
				// Copy entry + inject rule.Set
				newEntry := make(map[string]interface{})
				for k, v := range entry {
					newEntry[k] = v
				}
				for k, v := range rule.Set {
					newEntry[k] = v
				}
				out.Items = append(out.Items, newEntry)
			}
		}

		out.SourceCount = len(inputFiles)

		// Final output object
		outJSON := map[string]interface{}{}
		for k, v := range top {
			outJSON[k] = v
		}

		outJSON["entry_count"] = len(out.Items)
		outJSON["generated_at"] = time.Now().Format(time.RFC3339)
		outJSON["source_count"] = out.SourceCount
		outJSON["items"] = out.Items

		if err := writeJSON(rule.Output, outJSON); err != nil {
			return fmt.Errorf("writing output %s: %w", rule.Output, err)
		}
	}

	return nil
}

func matchEntry(entry map[string]interface{}, match map[string]interface{}) bool {
	if len(match) == 0 {
		return true
	}
	for k, v := range match {
		if entry[k] != v {
			return false
		}
	}
	return true
}

func writeJSON(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
