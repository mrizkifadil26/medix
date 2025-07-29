package enricher

import (
	"encoding/json"
	"flag"
	"os"
)

type CLIArgs struct {
	ConfigPath string
	Input      string
	Output     string
	Kind       string
}

func ParseCLI() (*Config, error) {
	var args CLIArgs
	flag.StringVar(&args.ConfigPath, "config", "", "Path to config JSON")

	flag.StringVar(&args.Input, "input", "", "Override: input file path")
	flag.StringVar(&args.Output, "output", "", "Override: output file path")
	flag.StringVar(&args.Kind, "kind", "", "Override: only kind (movie|tv)")
	flag.Parse()

	// Start with config file if provided
	cfg := &Config{}
	if args.ConfigPath != "" {
		data, err := os.ReadFile(args.ConfigPath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Apply CLI overrides
	if args.Input != "" {
		cfg.InputFile = args.Input
	}
	if args.Output != "" {
		cfg.OutputFile = args.Output
	}
	if args.Kind != "" {
		cfg.OnlyKind = args.Kind
	}

	return cfg, nil
}
