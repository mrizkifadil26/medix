package normalizer

import (
	"flag"
	"fmt"
)

type CLIArgs struct {
	ConfigPath *string
	Config     Config
}

func ParseCLI() (*CLIArgs, error) {
	var (
		configPath = flag.String("config", "", "Path to config file (JSON or YAML)")
		outputPath = flag.String("output", "", "Output result path")
		root       = flag.String("root", "", "Root directory to scan")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)

	flag.Parse()

	// Ensure config file is provided
	if configPath == nil || *configPath == "" {
		return nil, fmt.Errorf("missing required -config argument")
	}

	var cfg Config
	shouldPopulate := false

	if root != nil && *root != "" {
		cfg.Root = *root
		shouldPopulate = true
	}

	if verbose != nil && *verbose {
		cfg.Verbose = *verbose
		shouldPopulate = true
	}

	if outputPath != nil && *outputPath != "" {
		cfg.OutputPath = *outputPath
		shouldPopulate = true
	}

	args := &CLIArgs{
		ConfigPath: configPath,
	}

	if shouldPopulate {
		args.Config = cfg
	}

	// fmt.Println(*args.ConfigPath)

	return args, nil

	// configFile, err := os.ReadFile(configAbsPath)
	// if err != nil {
	// 	return CLIArgs{}, fmt.Errorf("failed to read config file: %w", err)
	// }

	// var config Config
	// if err := json.Unmarshal(configFile, &config); err != nil {
	// 	return CLIArgs{}, fmt.Errorf("failed to parse config JSON: %w", err)
	// }

	// // Determine final input path: CLI > config
	// // Determine input path: CLI overrides config
	// finalInput := config.File
	// if inputPath != "" {
	// 	finalInput = inputPath
	// }

	// if finalInput == "" {
	// 	return CLIArgs{}, fmt.Errorf("input path not provided (either via --input or config.file)")
	// }

	// return CLIArgs{
	// 	Input:      finalInput,
	// 	OutputPath: outputPath,
	// 	ConfigPath: configPath,
	// 	Config:     config,
	// }, nil
}
