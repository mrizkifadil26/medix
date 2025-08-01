package scannerV2

type Config struct {
	Root    string  `json:"root" yaml:"root"`
	Options Options `json:"options" yaml:"options"`
}

type Options struct {
	Mode        string   `json:"mode" yaml:"mode"`                                   // REQUIRED: controls behavior, never omit
	Exts        []string `json:"exts,omitempty" yaml:"exts,omitempty"`               // OPTIONAL: filters; empty = all
	Exclude     []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`         // OPTIONAL: path filtering
	Depth       int      `json:"depth" yaml:"depth"`                                 // REQUIRED: traversal logic
	OnlyLeaf    bool     `json:"onlyLeaf,omitempty" yaml:"onlyLeaf,omitempty"`       // OPTIONAL: feature toggle
	LeafDepth   int      `json:"leafDepth,omitempty" yaml:"leafDepth,omitempty"`     // OPTIONAL: advanced control
	SkipEmpty   bool     `json:"skipEmpty,omitempty" yaml:"skipEmpty,omitempty"`     // OPTIONAL: cosmetic/efficiency
	Concurrency int      `json:"concurrency,omitempty" yaml:"concurrency,omitempty"` // OPTIONAL: perf tuning
	Verbose     bool     `json:"verbose,omitempty" yaml:"verbose,omitempty"`         // OPTIONAL: cosmetic/debug

	// NEW — for subentry scanning
	SubEntries SubentriesMode `json:"subEntries"`        // If true, collect files inside dirs
	SubDepth   int            `json:"subDepth"`          // NEW: limit subentry depth, -1 = unlimited, 0 = none
	SubExts    []string       `json:"subExts,omitempty"` // Optional: file filters for subentries
}

func (c *Config) ApplyDefaults() {
	if c.Options.Mode == "" {
		c.Options.Mode = "files"
	}

	if c.Options.Depth == 0 {
		c.Options.Depth = 1
	}
	// All other fields (Exclude, Exts, OnlyLeaf, LeafDepth, SkipEmpty, etc.)
	// are optional and left as-is to respect user intent.
}
