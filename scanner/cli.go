package scanner

import (
	"flag"
	"fmt"
)

type CLIArgs struct {
	ConfigPath *string
	Config     Config // Full config with all options
}

func ParseCLI() (*CLIArgs, error) {
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
		fmt.Println(helpText)
	}

	var (
		configPath = flag.String("config", "", "Path to config file (JSON or YAML)")
		outputPath = flag.String("output", "", "Output result path")
		root       = flag.String("root", "", "Root directory to scan")
		mode       = flag.String("mode", "", "Scan mode: files, dirs, or mixed")
		depth      = flag.Int("depth", -1, "Max scan depth")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)

	flag.Parse()

	if configPath == nil || *configPath == "" {
		return nil, fmt.Errorf("missing required -config argument")
	}

	var cfg Config
	shouldPopulate := false

	if root != nil && *root != "" {
		cfg.Root = *root
		shouldPopulate = true
	}
	if mode != nil && *mode != "" {
		cfg.Options.Mode = *mode
		shouldPopulate = true
	}
	if depth != nil && *depth != -1 {
		cfg.Options.Depth = *depth
		shouldPopulate = true
	}
	if verbose != nil && *verbose {
		cfg.Options.Verbose = *verbose
		shouldPopulate = true
	}
	if outputPath != nil && *outputPath != "" {
		if cfg.Output == nil {
			cfg.Output = &OutputOptions{}
		}

		cfg.Output.OutputPath = *outputPath
		shouldPopulate = true
	}

	args := &CLIArgs{
		ConfigPath: configPath,
	}

	if shouldPopulate {
		args.Config = cfg
	}

	return args, nil
}
