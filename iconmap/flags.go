package iconmap

import (
	"flag"
	"fmt"
	"os"
)

type CLIFlags struct {
	ConfigPath string

	Source string
	Type   string // Should be "movies" or "tv"
	Name   string
	Output string
}

func ParseFlags() CLIFlags {
	var f CLIFlags

	flag.StringVar(&f.ConfigPath, "config", "", "Path to config file (JSON). If provided, inline flags are ignored.")

	flag.StringVar(&f.Source, "source", "", "Path to source directory (inline)")
	flag.StringVar(&f.Type, "type", "", `Media type: "movies" or "tv" (inline)`)
	flag.StringVar(&f.Name, "name", "", "Name of the source (inline)")
	flag.StringVar(&f.Output, "output", "", "Path to output JSON file (inline)")

	flag.Parse()

	// Validate
	if f.ConfigPath == "" {
		// No config file, using inline mode — validate required fields
		if f.Source == "" || f.Type == "" || f.Name == "" || f.Output == "" {
			fmt.Fprintln(os.Stderr, "❌ Missing required inline flags: --source, --type, --name, and --output are all required if --config is not provided.")
			flag.Usage()
			os.Exit(1)
		}

		if f.Type != "movies" && f.Type != "tv" {
			fmt.Fprintln(os.Stderr, `❌ Invalid type: must be either "movies" or "tv".`)
			os.Exit(1)
		}
	}

	return f
}
