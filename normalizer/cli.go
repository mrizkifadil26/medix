package normalizer

import (
	"flag"
	"strings"
)

type CLIArgs struct {
	Input      string
	OutputPath string
	ConfigPath string
	Config     Config
}

func ParseCLI() CLIArgs {
	var args CLIArgs
	config := &args.Config
	var stepsStr string

	flag.StringVar(&args.Input, "input", "", "Input string, JSON array, or path to scan JSON")
	flag.StringVar(&args.OutputPath, "output", "", "Output file (optional)")
	flag.StringVar(&args.ConfigPath, "config", "", "Path to config JSON (optional)")
	flag.StringVar(&stepsStr, "steps", "", "Comma-separated normalization steps")
	flag.Parse()

	config.Steps = splitAndTrim(stepsStr)

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
