package scannerV2

import (
	"flag"
)

type CLIArgs struct {
	ConfigPath *string

	// Partial overrides (nullable)
	OutputPath *string
	Root       *string
	Mode       *string
	Depth      *int
	Verbose    *bool

	Config Config // Full config with all options
}

func ParseCLI() *CLIArgs {
	flag.Usage = func() {
		helpText := `
Usage: scan [OPTIONS]

Options:
  -config      Path to config file (JSON or YAML)
  -output      Path to output result (optional)
  -root        Root directory to scan (required if no config file)
  -mode        Scan mode: files, dirs, or mixed (default: files)
  -depth       Max scan depth (0 = top-level, -1 = unlimited) (default: 1)
  -verbose     Enable verbose logging

Example:
  scan -root ./media -mode dirs -depth 2 -verbose

If -config is provided, it overrides everything except -output.
`
		println(helpText)
	}

	// var args CLIArgs
	// options := &args.Config.Options
	// var extStr, excludeStr string
	var (
		configPath = flag.String("config", "", "Path to config file (JSON or YAML)")
		outputPath = flag.String("output", "", "Output result path")
		root       = flag.String("root", "", "Root directory to scan")
		mode       = flag.String("mode", "", "Scan mode: files, dirs, or mixed")
		depth      = flag.Int("depth", -9999, "Max scan depth")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)

	flag.Parse()

	// Global CLI-only flags
	// flag.StringVar(&configPath, "config", "", "Path to config file (JSON or YAML)")
	// flag.StringVar(&outputPath, "output", "", "Path to output result (optional)")
	// flag.StringVar(&root, "root", "", "Root directory to scan (required if no config file)")

	// Config: scan options
	// flag.StringVar(&options.Mode, "mode", "files", "Scan mode: files or dirs")
	// flag.StringVar(&extStr, "exts", ".mkv,.mp4", "Comma-separated file extensions")
	// flag.StringVar(&excludeStr, "exclude", "", "Comma-separated excluded paths")

	// flag.IntVar(&options.Depth, "depth", 1, "Max scan depth (0 for top-level, -1 for unlimited)")
	// flag.IntVar(&options.LeafDepth, "leaf-depth", 0, "Min leaf directory depth")
	// flag.IntVar(&options.Concurrency, "concurrency", 0, "Number of concurrent workers (0 = auto)")

	// flag.BoolVar(&options.OnlyLeaf, "only-leaf", false, "Only scan leaf directories")
	// flag.BoolVar(&options.SkipEmpty, "skip-empty", false, "Skip empty directories")
	// flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()

	// Convert comma-separated strings into slices
	// options.Exts = splitAndTrim(extStr)
	// options.Exclude = splitAndTrim(excludeStr)

	cfg := Config{}
	// Only assign if explicitly set (non-empty or non-placeholder)
	if root != nil && *root != "" {
		cfg.Root = root
	}

	if outputPath != nil && *outputPath != "" {
		cfg.Output.OutputPath = outputPath
	}

	if verbose != nil && flag.Lookup("verbose").Value.String() == "true" {
		cfg.Verbose = verbose
	}

	if mode != nil && *mode != "" {
		cfg.Options.Mode = *mode
	}

	if depth != nil && *depth != -9999 {
		cfg.Options.Depth = *depth
	}

	return &CLIArgs{
		ConfigPath: configPath,
		OutputPath: outputPath,
		Config:     cfg,
	}
}

// func splitAndTrim(s string) []string {
// 	if s == "" {
// 		return nil
// 	}
// 	parts := strings.Split(s, ",")
// 	for i, p := range parts {
// 		parts[i] = strings.TrimSpace(p)
// 	}
// 	return parts
// }
