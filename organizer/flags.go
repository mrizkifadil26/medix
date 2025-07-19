package organizer

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mrizkifadil26/medix/utils"
)

type Flags struct {
	ConfigPath string

	Mode      string // Should be "movies" or "tv"
	InputPath string
	Sources   utils.ArrayFlags
	Output    string
}

func Parse() Flags {
	var f Flags

	flag.StringVar(&f.ConfigPath, "config", "", "Path to config file (JSON). If provided, inline flags are ignored.")

	flag.StringVar(&f.Mode, "mode", "", `Operation mode: "preview" or "apply" (inline)`)
	flag.StringVar(&f.InputPath, "json", "", "Path to input JSON file (inline)")
	flag.Var(&f.Sources, "sources", "Path(s) to source directory (comma-separated or repeated)")
	flag.StringVar(&f.Output, "output", "", "Path to output JSON file (inline)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
Usage:
  organize --mode=preview --json=tv.json --sources=Icons --output=out.json
  organize --mode=apply --json=tv.json --sources=Icons,OtherIcons --output=result.json
  organize --config=config.json

Flags:`)
		flag.PrintDefaults()
	}

	flag.Parse()

	// After flag.Parse()
	if f.ConfigPath != "" {
		var cfg Config
		utils.LoadJSON(f.ConfigPath, &cfg)

		f.Mode = cfg.Mode
		f.InputPath = cfg.MediaInput
		f.Output = cfg.TargetDir
		f.Sources = cfg.IconSources
	}

	if err := f.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "❌", err)
		flag.Usage()
		os.Exit(1)
	}

	return f
}

func (f Flags) Validate() error {
	if f.ConfigPath == "" {
		if f.Mode == "" || f.InputPath == "" || len(f.Sources) == 0 || f.Output == "" {
			return errors.New("missing required inline flags")
		}
	} else {
		// Validate loaded config fields
		if f.Mode != "preview" && f.Mode != "apply" {
			return fmt.Errorf("invalid mode in config: must be 'preview' or 'apply'")
		}
	}
	return nil
}
