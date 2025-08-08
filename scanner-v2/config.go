package scannerV2

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils"
)

type Config struct {
	Root    *string   `json:"root" yaml:"root"`
	Tags    *[]string `json:"tags,omitempty" yaml:"tags,omitempty"` // Optional job/context tags
	Verbose *bool     `json:"verbose" yaml:"verbose"`

	Options *ScanOptions   `json:"options" yaml:"options"`
	Output  *OutputOptions `json:"output" yaml:"output"` // Output format and options
	// FileRules []Rule      `json:"fileRules,omitempty"`
	// DirRules  []Rule      `json:"dirRules,omitempty"`
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
	IncludeChildren  bool   `json:"includeChildren,omitempty" yaml:"includeChildren,omitempty"`
	OnlyLeaf         bool   `json:"onlyLeaf,omitempty" yaml:"onlyLeaf,omitempty"` // OPTIONAL: feature toggle
	Trace            bool   `json:"trace,omitempty" yaml:"trace,omitempty"`
}

type OutputOptions struct {
	Format     *string `json:"format" yaml:"format"`                             // Output format: "json", "yaml", etc
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

// config.go
func (baseConfig *Config) ApplyDefaults() *Config {
	defaultCfg := DefaultConfig()

	merged, err := utils.MergeDeep(*baseConfig, defaultCfg)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return &merged
}

func (cfg *Config) PrettyPrint() {
	fmt.Println("📦 Scanning Configuration")

	// Root-level fields
	if cfg.Root != nil {
		printRow("Root", *cfg.Root, "Root directory to scan")
	}
	if cfg.Tags != nil {
		printRow("Tags", fmt.Sprintf("%v", *cfg.Tags), "Optional tags for context")
	}
	if cfg.Verbose != nil {
		printRow("Verbose", fmt.Sprintf("%v", *cfg.Verbose), "Verbose logging enabled")
	}

	// Options section
	if cfg.Options != nil {
		fmt.Println()
		fmt.Println("🛠️  Scan Options:")
		printRow("Mode", cfg.Options.Mode, "Scan mode: files, dirs, or mixed")
		printRow("Depth", fmt.Sprintf("%d", cfg.Options.Depth), "How deep to traverse directories")
		printRow("SkipEmpty", fmt.Sprintf("%v", cfg.Options.SkipEmpty), "Skip empty directories")
		printRow("IncludeRootFiles", fmt.Sprintf("%v", cfg.Options.IncludeRootFiles), "Include root-level files")
		printRow("IncludeChildren", fmt.Sprintf("%v", cfg.Options.IncludeChildren), "Include children dirs and files")
		printRow("OnlyLeaf", fmt.Sprintf("%v", cfg.Options.OnlyLeaf), "Only include leaf-level directories")
	}

	// Output section
	if cfg.Output != nil {
		fmt.Println()
		fmt.Println("📤 Output Options:")
		if cfg.Output.Format != nil {
			printRow("Format", *cfg.Output.Format, "Output format: json, yaml, etc.")
		}

		if cfg.Output.OutputPath != nil {
			printRow("OutputPath", *cfg.Output.OutputPath, "Path to save output (optional)")
		}

		printRow("IncludeErrors", fmt.Sprintf("%v", cfg.Output.IncludeErrors), "Include error info in output")
		printRow("IncludeWarnings", fmt.Sprintf("%v", cfg.Output.IncludeWarnings), "Include warnings in output")
		printRow("IncludeStats", fmt.Sprintf("%v", cfg.Output.IncludeStats), "Include detailed scan stats")
	}
}

func printRow(key, value, comment string) {
	const keyWidth = 18
	const valWidth = 24

	keyStr := fmt.Sprintf("%-*s", keyWidth, key)
	valStr := fmt.Sprintf("%-*s", valWidth, ":"+value)
	fmt.Printf("%s %s # %s\n", keyStr, valStr, comment)
}
