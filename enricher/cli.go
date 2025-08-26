package enricher

import (
	"flag"
	"fmt"
)

type CLIArgs struct {
	ConfigPath *string
	Refresh    *bool
	Config     Config
}

func ParseCLI() (*CLIArgs, error) {
	var (
		configPath = flag.String("config", "", "Path to config file (JSON or YAML)")
		outputPath = flag.String("output", "", "Output result path")
		root       = flag.String("root", "", "Root directory to scan")
		refresh    = flag.Bool("refresh", false, "Force enrichment by ignoring cache")
	)

	flag.Parse()

	// Ensure config file is provided
	if configPath == nil || *configPath == "" {
		return nil, fmt.Errorf("missing required -config argument")
	}

	// Start with config file if provided
	var cfg Config
	shouldPopulate := false

	if root != nil && *root != "" {
		cfg.Root = *root
		shouldPopulate = true
	}

	if outputPath != nil && *outputPath != "" {
		cfg.Output = *outputPath
		shouldPopulate = true
	}

	args := &CLIArgs{
		ConfigPath: configPath,
		Refresh:    refresh,
	}

	if shouldPopulate {
		args.Config = cfg
	}

	return args, nil
}
