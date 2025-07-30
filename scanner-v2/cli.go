package scannerV2

import (
	"flag"
	"strings"
)

type CLIArgs struct {
	ConfigPath string
	OutputPath string
	ScanConfig
}

func ParseCLI() CLIArgs {
	var cli CLIArgs
	var exts, exclude string

	flag.StringVar(&cli.Root, "root", "", "Root directory to scan (required)")
	flag.StringVar(&cli.Mode, "mode", "files", "Scan mode: files or dirs")
	flag.IntVar(&cli.Depth, "depth", -1, "Max depth for directory traversal (-1 for unlimited)")
	flag.StringVar(&exts, "exts", ".mkv,.mp4", "Comma-separated file extensions (files mode only)")
	flag.StringVar(&exclude, "exclude", "", "Comma-separated list of excluded paths")
	flag.BoolVar(&cli.OnlyLeaf, "only-leaf", false, "Only scan leaf directories (no subfolders)")
	flag.BoolVar(&cli.SkipEmpty, "skip-empty", false, "Skip empty directories from scan output")
	flag.BoolVar(&cli.Verbose, "verbose", false, "Enable verbose logging")

	flag.StringVar(&cli.ConfigPath, "config", "", "Path to JSON config file")
	flag.StringVar(&cli.OutputPath, "output", "", "Path to write JSON output (if empty, print to stdout)")

	flag.Parse()

	cli.Exts = splitAndTrim(exts)
	cli.Exclude = splitAndTrim(exclude)

	return cli
}

// func (cli CLIArgs) OverrideConfig(cfg *ScanConfig) {
// 	if cli.Root != "" {
// 		cfg.Root = cli.Root
// 	}
// 	if cli.Mode != "" {
// 		cfg.Mode = cli.Mode
// 	}
// 	if len(cli.Exts) > 0 {
// 		cfg.Exts = cli.Exts
// 	}
// 	if cli.Depth >= 0 {
// 		cfg.Depth = cli.Depth
// 	}
// 	if len(cli.Exclude) > 0 {
// 		cfg.Exclude = cli.Exclude
// 	}
// 	if cli.OnlyLeaf {
// 		cfg.OnlyLeaf = true
// 	}
// 	if cli.LeafDepth > 0 {
// 		cfg.LeafDepth = cli.LeafDepth
// 	}
// 	if cli.SkipEmpty {
// 		cfg.SkipEmpty = true
// 	}
// 	if cli.Verbose {
// 		cfg.Verbose = true
// 	}
// }

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
