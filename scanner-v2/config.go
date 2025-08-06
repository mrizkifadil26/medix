package scannerV2

import "github.com/mrizkifadil26/medix/utils"

type Config struct {
	Root    *string  `json:"root" yaml:"root"`
	Tags    []string `json:"tags,omitempty" yaml:"tags,omitempty"` // Optional job/context tags
	Verbose *bool    `json:"verbose" yaml:"verbose"`

	Options *ScanOptions `json:"options" yaml:"options"`
	// FileRules []Rule      `json:"fileRules,omitempty"`
	// DirRules  []Rule      `json:"dirRules,omitempty"`

	Output *OutputOptions `json:"output" yaml:"output"` // Output format and options
}

// type Options struct {
// 	Mode        string   `json:"mode" yaml:"mode"`                                   // REQUIRED: controls behavior, never omit
// 	Exts        []string `json:"exts,omitempty" yaml:"exts,omitempty"`               // OPTIONAL: filters; empty = all
// 	Exclude     []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`         // OPTIONAL: path filtering
// 	Depth       int      `json:"depth" yaml:"depth"`                                 // REQUIRED: traversal logic
// 	OnlyLeaf    bool     `json:"onlyLeaf,omitempty" yaml:"onlyLeaf,omitempty"`       // OPTIONAL: feature toggle
// 	LeafDepth   int      `json:"leafDepth,omitempty" yaml:"leafDepth,omitempty"`     // OPTIONAL: advanced control
// 	SkipEmpty   bool     `json:"skipEmpty,omitempty" yaml:"skipEmpty,omitempty"`     // OPTIONAL: cosmetic/efficiency
// 	Concurrency int      `json:"concurrency,omitempty" yaml:"concurrency,omitempty"` // OPTIONAL: perf tuning

// 	// NEW — for subentry scanning
// 	SubEntries SubentriesMode `json:"subEntries"`        // If true, collect files inside dirs
// 	SubDepth   int            `json:"subDepth"`          // NEW: limit subentry depth, -1 = unlimited, 0 = none
// 	SubExts    []string       `json:"subExts,omitempty"` // Optional: file filters for subentries
// }

type ScanOptions struct {
	Mode             string `json:"mode" yaml:"mode"`                               // "files", "dirs", "mixed"
	Depth            int    `json:"depth" yaml:"depth"`                             // REQUIRED: traversal logic
	SkipEmpty        bool   `json:"skipEmpty,omitempty" yaml:"skipEmpty,omitempty"` // OPTIONAL: skip empty directories
	IncludeRootFiles bool   `json:"includeRootFiles,omitempty" yaml:"includeRootFiles,omitempty"`
	OnlyLeaf         bool   `json:"onlyLeaf,omitempty" yaml:"onlyLeaf,omitempty"` // OPTIONAL: feature toggle
}

type OutputOptions struct {
	Format     string  `json:"format" yaml:"format"`                             // Output format: "json", "yaml", etc
	OutputPath *string `json:"outputPath,omitempty" yaml:"outputPath,omitempty"` // Optional output file path

	IncludeErrors   bool `json:"includeErrors,omitempty" yaml:"includeErrors,omitempty"`     // Include errors in output
	IncludeWarnings bool `json:"includeWarnings,omitempty" yaml:"includeWarnings,omitempty"` // Include warnings in output
	IncludeStats    bool `json:"includeStats,omitempty" yaml:"includeStats,omitempty"`       // Include detailed stats
}

type Rule struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions,omitempty"` // For files
	Patterns   []string `json:"patterns,omitempty"`   // Glob or regex
	MinFiles   int      `json:"minFiles,omitempty"`   // For dir rule
	MaxFiles   int      `json:"maxFiles,omitempty"`
}

// func (c *Config) ApplyDefaults() {
// 	if c.Options.Mode == "" {
// 		c.Options.Mode = "files"
// 	}

// 	if c.Options.Depth == 0 {
// 		c.Options.Depth = 1
// 	}

// 	// All other fields (Exclude, Exts, OnlyLeaf, LeafDepth, SkipEmpty, etc.)
// 	// are optional and left as-is to respect user intent.
// }

// config.go
func (c *Config) ApplyDefaults() {
	def := DefaultConfig()

	utils.MergeDeep(c, &def)
}
