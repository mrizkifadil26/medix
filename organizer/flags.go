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

	Mode      string // preview | apply (must be CLI-provided)
	InputPath string
	Sources   utils.KVFlags
	TargetDir string

	ReportPath string // Where to save/read the changes.json
	ExportPath string // Only used in preview (optional)

	IconSources []IconSource // internally resolved
}

func Parse() Flags {
	var f Flags

	flag.StringVar(&f.ConfigPath, "config", "", "Path to config file (JSON). If provided, inline flags are ignored.")
	flag.StringVar(&f.Mode, "mode", "", `Operation mode: "preview" or "apply"`)

	flag.StringVar(&f.InputPath, "json", "", "Path to input JSON file (usually changes report)")
	flag.Var(&f.Sources, "sources", "Icon sources in label:path format or just path")
	flag.StringVar(&f.TargetDir, "target", "", "Target directory for organizing (only in apply mode)")
	flag.StringVar(&f.ReportPath, "report", "", "Path to the main changes report file (default: data/organize/changes.json)")
	flag.StringVar(&f.ExportPath, "export", "", "Optional: also export preview report to another location")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
Usage (choose one mode):

  # Inline CLI
  organize --mode=preview --json=media-index.json --sources=Icons --report=data/organize/changes.json
  organize --mode=apply   --json=data/organize/changes.json --target=Media/Movies

  # Config-based (with --mode override)
  organize --config=config.json --mode=preview

Flags:`)
		flag.PrintDefaults()
	}

	flag.Parse()

	// If config is used, parse JSON config file
	if f.ConfigPath != "" {
		var cfg Config
		utils.LoadJSON(f.ConfigPath, &cfg)

		f.InputPath = cfg.MediaInput
		f.TargetDir = cfg.TargetDir
		f.IconSources = cfg.IconSources
		f.ReportPath = cfg.ReportPath
	}

	// Set default report path if not specified
	if f.ReportPath == "" {
		f.ReportPath = "data/organize/changes.json"
	}

	// Convert sources into []IconSource if not using config
	if f.ConfigPath == "" {
		for source, path := range f.Sources.ToMap() {
			f.IconSources = append(f.IconSources, IconSource{
				Source: source,
				Path:   path,
			})
		}
	}

	// ✅ Validate after parsing
	if err := f.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "❌", err)
		flag.Usage()
		os.Exit(1)
	}

	return f
}

func (f Flags) Validate() error {
	if f.Mode != "preview" && f.Mode != "apply" {
		return errors.New("missing or invalid --mode (must be 'preview' or 'apply')")
	}

	if f.InputPath == "" {
		return errors.New("missing input --json (media index or changes report)")
	}

	if f.Mode == "preview" && len(f.IconSources) == 0 {
		return errors.New("missing --sources in preview mode (icon sources)")
	}

	if f.Mode == "apply" && f.TargetDir == "" {
		return errors.New("missing --target in apply mode")
	}

	return nil
}
