package scanner

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils"
)

type Config struct {
	Root string   `json:"root" yaml:"root"`
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"` // Optional job/context tags

	Options *ScanOptions   `json:"options,omitempty" yaml:"options,omitempty"`
	Output  *OutputOptions `json:"output,omitempty" yaml:"output,omitempty"` // Output format and options
}

type ScanOptions struct {
	Mode             string `json:"mode" yaml:"mode"`                               // "files", "dirs", "mixed"
	Depth            int    `json:"depth" yaml:"depth"`                             // REQUIRED: traversal logic
	SkipEmpty        bool   `json:"skipEmpty,omitempty" yaml:"skipEmpty,omitempty"` // OPTIONAL: skip empty directories
	SkipRoot         bool   `json:"skipRoot,omitempty" yaml:"skipRoot,omitempty"`   // OPTIONAL: skip root directory
	MinIncludeDepth  int    `json:"minIncludeDepth,omitempty" yaml:"minIncludeDepth,omitempty"`
	IncludeHidden    bool   `json:"includeHidden,omitempty" yaml:"includeHidden,omitempty"` // Include hidden files/dirs
	IncludeRootFiles bool   `json:"includeRootFiles,omitempty" yaml:"includeRootFiles,omitempty"`
	IncludeChildren  bool   `json:"includeChildren,omitempty" yaml:"includeChildren,omitempty"`
	OnlyLeaf         bool   `json:"onlyLeaf,omitempty" yaml:"onlyLeaf,omitempty"` // OPTIONAL: feature toggle

	Verbose bool `json:"verbose,omitempty" yaml:"verbose,omitempty"`
	Debug   bool `json:"debug,omitempty" yaml:"debug,omitempty"` // Enable DEBUG logging
	Trace   bool `json:"trace,omitempty" yaml:"trace,omitempty"` // Enable TRACE logging

	IncludePatterns []string `json:"includePatterns,omitempty" yaml:"includePatterns,omitempty"` // Glob patterns to include
	ExcludePatterns []string `json:"excludePatterns,omitempty" yaml:"excludePatterns,omitempty"` // Glob patterns to exclude
	IncludeExts     []string `json:"includeExts,omitempty" yaml:"includeExts,omitempty"`         // File extensions to include
	ExcludeExts     []string `json:"excludeExts,omitempty" yaml:"excludeExts,omitempty"`         // File extensions to exclude

	EnableProgress bool `json:"enableProgress,omitempty" yaml:"enableProgress,omitempty"` // Show real-time progress during scan

	StopOnError bool `json:"stopOnError,omitempty" yaml:"stopOnError,omitempty"` // Stop walking on first error
	SkipOnError bool `json:"skipOnError,omitempty" yaml:"skipOnError,omitempty"` // Skip entries that cause errors
	Concurrency int  `json:"concurrency,omitempty" yaml:"concurrency,omitempty"` // Concurrency level for processing
}

type OutputOptions struct {
	Format     string `json:"format,omitempty" yaml:"format,omitempty"`         // Output format: "json", "yaml", etc
	OutputPath string `json:"outputPath,omitempty" yaml:"outputPath,omitempty"` // Optional output file path

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

	section := func(title string) { fmt.Printf("\n%s\n", title) }

	// Root-level fields
	if cfg.Root != "" {
		printRow("Root", cfg.Root, "Root directory to scan")
	}

	if cfg.Tags != nil {
		printRow("Tags", fmt.Sprintf("%v", cfg.Tags), "Optional tags for context")
	}

	// Options
	if cfg.Options != nil {
		section("🛠️  Scan Options:")
		rows := []struct {
			Name, Value, Desc string
		}{
			{"Mode", cfg.Options.Mode, "Scan mode: files, dirs, or mixed"},
			{"Depth", fmt.Sprintf("%d", cfg.Options.Depth), "Traversal depth (0=root only)"},

			{"Log Level", cfg.logLevelString(), "Logging verbosity level (verbose, debug, trace)"},

			{"OnlyLeaf", fmt.Sprintf("%v", cfg.Options.OnlyLeaf), "Only leaf directories"},
			{"SkipEmpty", fmt.Sprintf("%v", cfg.Options.SkipEmpty), "Skip empty directories"},
			{"SkipRoot", fmt.Sprintf("%v", cfg.Options.SkipRoot), "Skip root directory"},
			{"IncludeHidden", fmt.Sprintf("%v", cfg.Options.IncludeHidden), "Include hidden files/dirs"},
			{"IncludeRootFiles", fmt.Sprintf("%v", cfg.Options.IncludeRootFiles), "Include root-level files"},
			{"IncludeChildren", fmt.Sprintf("%v", cfg.Options.IncludeChildren), "Include child dirs/files"},
			{"MinIncludeDepth", fmt.Sprintf("%v", cfg.Options.MinIncludeDepth), "Include child dirs/files"},
			{"Trace", fmt.Sprintf("%v", cfg.Options.Trace), "Enable tracing"},
			{"EnableProgress", fmt.Sprintf("%v", cfg.Options.EnableProgress), "Show progress updates"},
			{"StopOnError", fmt.Sprintf("%v", cfg.Options.StopOnError), "Stop on first error"},
			{"SkipOnError", fmt.Sprintf("%v", cfg.Options.SkipOnError), "Skip entries on error"},
			{"Concurrency", fmt.Sprintf("%d", cfg.Options.Concurrency), "Worker concurrency"},
		}
		for _, r := range rows {
			printRow(r.Name, r.Value, r.Desc)
		}
	}

	// Output section
	if cfg.Output != nil {
		section("📤 Output Options:")

		printRow("Format", cfg.Output.Format, "Output format: json, yaml, etc.")
		printRow("OutputPath", cfg.Output.OutputPath, "Path to save output (optional)")
		printRow("IncludeErrors", fmt.Sprintf("%v", cfg.Output.IncludeErrors), "Include error info in output")
		printRow("IncludeWarnings", fmt.Sprintf("%v", cfg.Output.IncludeWarnings), "Include warnings in output")
		printRow("IncludeStats", fmt.Sprintf("%v", cfg.Output.IncludeStats), "Include detailed scan stats")
	}

	fmt.Println()
}

func printRow(key, value, comment string) {
	const keyWidth = 18
	const valWidth = 24

	keyStr := fmt.Sprintf("%-*s", keyWidth, key)
	valStr := fmt.Sprintf("%-*s", valWidth, ":"+value)
	fmt.Printf("%s %s # %s\n", keyStr, valStr, comment)
}

func (cfg *Config) logLevelString() string {
	opts := cfg.Options

	if opts.Verbose {
		return "verbose"
	}

	if opts.Debug {
		return "debug"
	}

	if opts.Trace {
		return "trace"
	}

	return "normal"
}
