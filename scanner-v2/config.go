package scannerV2

import (
	"encoding/json"
	"fmt"
	"os"
)

type ScanConfig struct {
	Root        string   `json:"root"`
	Mode        string   `json:"mode"` // "files" or "dirs"
	Exts        []string `json:"exts"`
	Depth       int      `json:"depth"`
	Exclude     []string `json:"exclude"`
	OnlyLeaf    bool     `json:"onlyLeaf"`
	LeafDepth   int      `json:"leafDepth"` // 0 = default, 1 = leaf-1, 2 = leaf-2, etc.
	SkipEmpty   bool     `json:"skipEmpty"`
	Concurrency int      `json:"concurrency"` // worker count
	Verbose     bool     `json:"verbose"`

	// NEW — for subentry scanning
	SubEntries SubentriesMode `json:"subEntries"`        // If true, collect files inside dirs
	SubDepth   int            `json:"subDepth"`          // NEW: limit subentry depth, -1 = unlimited, 0 = none
	SubExts    []string       `json:"subExts,omitempty"` // Optional: file filters for subentries
}

func (c ScanConfig) ToOptions() ScanOptions {
	return ScanOptions{
		Mode:        c.Mode,
		Exts:        c.Exts,
		Depth:       c.Depth,
		Exclude:     c.Exclude,
		OnlyLeaf:    c.OnlyLeaf,
		LeafDepth:   c.LeafDepth,
		SkipEmpty:   c.SkipEmpty,
		Verbose:     c.Verbose,
		Concurrency: c.Concurrency,

		SubEntries: c.SubEntries,
		SubExts:    c.SubExts,
		SubDepth:   c.SubDepth,
	}
}

func LoadConfig(path string) (*ScanConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg ScanConfig
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &cfg, nil
}
