package normalizer

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type CLIArgs struct {
	Input      string
	OutputPath string
	ConfigPath string
	Config     Config
}

func ParseCLI() (CLIArgs, error) {
	var (
		inputPath  string
		outputPath string
		configPath string
	)

	// Define CLI flags
	flag.StringVar(&inputPath, "input", "", "Override input file path (optional, defaults to 'file' in config)")
	flag.StringVar(&outputPath, "output", "", "Optional output file path")
	flag.StringVar(&configPath, "config", "", "Path to normalization config JSON (required)")
	flag.Parse()

	// Ensure config file is provided
	if configPath == "" {
		return CLIArgs{}, fmt.Errorf("missing required --config flag")
	}

	// Load config from file
	configAbsPath, err := filepath.Abs(configPath)
	if err != nil {
		return CLIArgs{}, fmt.Errorf("invalid config path: %w", err)
	}

	configFile, err := os.ReadFile(configAbsPath)
	if err != nil {
		return CLIArgs{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return CLIArgs{}, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Determine final input path: CLI > config
	// Determine input path: CLI overrides config
	finalInput := config.File
	if inputPath != "" {
		finalInput = inputPath
	}

	if finalInput == "" {
		return CLIArgs{}, fmt.Errorf("input path not provided (either via --input or config.file)")
	}

	return CLIArgs{
		Input:      finalInput,
		OutputPath: outputPath,
		ConfigPath: configPath,
		Config:     config,
	}, nil
}
