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

type ScanOptions struct {
	Mode             string `json:"mode" yaml:"mode"`                                       // "files", "dirs", "mixed"
	Depth            int    `json:"depth" yaml:"depth"`                                     // REQUIRED: traversal logic
	SkipEmpty        bool   `json:"skipEmpty,omitempty" yaml:"skipEmpty,omitempty"`         // OPTIONAL: skip empty directories
	IncludeHidden    bool   `json:"includeHidden,omitempty" yaml:"includeHidden,omitempty"` // Include hidden files/dirs
	IncludeRootFiles bool   `json:"includeRootFiles,omitempty" yaml:"includeRootFiles,omitempty"`
	IncludeChildren  bool   `json:"includeChildren,omitempty" yaml:"includeChildren,omitempty"`
	OnlyLeaf         bool   `json:"onlyLeaf,omitempty" yaml:"onlyLeaf,omitempty"` // OPTIONAL: feature toggle
	Trace            bool   `json:"trace,omitempty" yaml:"trace,omitempty"`

	EnableProgress bool `json:"enableProgress,omitempty" yaml:"enableProgress,omitempty"` // Show real-time progress during scan

	IncludePatterns []string `json:"includePatterns,omitempty" yaml:"includePatterns,omitempty"` // Glob patterns to include
	ExcludePatterns []string `json:"excludePatterns,omitempty" yaml:"excludePatterns,omitempty"` // Glob patterns to exclude
	IncludeExts     []string `json:"includeExts,omitempty" yaml:"includeExts,omitempty"`         // File extensions to include
	ExcludeExts     []string `json:"excludeExts,omitempty" yaml:"excludeExts,omitempty"`         // File extensions to exclude

	StopOnError bool `json:"stopOnError,omitempty" yaml:"stopOnError,omitempty"` // Stop walking on first error
	SkipOnError bool `json:"skipOnError,omitempty" yaml:"skipOnError,omitempty"` // Skip entries that cause errors
	Concurrency int  `json:"concurrency,omitempty" yaml:"concurrency,omitempty"` // Concurrency level for processing
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
func (c *Config) ApplyDefaults() error {
	defaults := DefaultConfig()

	return utils.MergeInto(c, &defaults, utils.MergeOptions{
		Overwrite: false,
		Recursive: true,
	})
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
