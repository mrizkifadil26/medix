package scannerV2

import (
	"flag"
	"strings"
)

type CLIArgs struct {
	ConfigPath string
	OutputPath string
	Config     Config
}

func ParseCLI() CLIArgs {
	flag.Usage = func() {
		helpText := `
Usage: scan [OPTIONS]

Options:
  -config         Path to config file (JSON or YAML)
  -output         Path to output result (optional)
  -root           Root directory to scan (required if no config file)

Scan Options:
  -mode           Scan mode: files or dirs (default: files)
  -exts           Comma-separated file extensions (default: .mkv,.mp4)
  -exclude        Comma-separated excluded paths
  -depth          Max scan depth (0 = top-level, -1 = unlimited) (default: 1)
  -leaf-depth     Minimum leaf directory depth
  -concurrency    Number of concurrent workers (0 = auto)
  -only-leaf      Only scan leaf directories
  -skip-empty     Skip empty directories
  -verbose        Enable verbose logging

Example:
  scan -root /media/movies -exts .mkv,.mp4 -depth 2 -only-leaf -verbose

If -config is provided, it overrides CLI flags (except -output).

`
		println(helpText)
	}

	var args CLIArgs
	options := &args.Config.Options
	var extStr, excludeStr string

	// Global CLI-only flags
	flag.StringVar(&args.ConfigPath, "config", "", "Path to config file (JSON or YAML)")
	flag.StringVar(&args.OutputPath, "output", "", "Path to output result (optional)")

	// Config: top-level
	flag.StringVar(&args.Config.Root, "root", "", "Root directory to scan (required if no config file)")

	// Config: scan options
	flag.StringVar(&options.Mode, "mode", "files", "Scan mode: files or dirs")
	flag.StringVar(&extStr, "exts", ".mkv,.mp4", "Comma-separated file extensions")
	flag.StringVar(&excludeStr, "exclude", "", "Comma-separated excluded paths")

	flag.IntVar(&options.Depth, "depth", 1, "Max scan depth (0 for top-level, -1 for unlimited)")
	flag.IntVar(&options.LeafDepth, "leaf-depth", 0, "Min leaf directory depth")
	flag.IntVar(&options.Concurrency, "concurrency", 0, "Number of concurrent workers (0 = auto)")

	flag.BoolVar(&options.OnlyLeaf, "only-leaf", false, "Only scan leaf directories")
	flag.BoolVar(&options.SkipEmpty, "skip-empty", false, "Skip empty directories")
	flag.BoolVar(&options.Verbose, "verbose", false, "Enable verbose logging")

	flag.Parse()

	// Convert comma-separated strings into slices
	options.Exts = splitAndTrim(extStr)
	options.Exclude = splitAndTrim(excludeStr)

	return args
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
