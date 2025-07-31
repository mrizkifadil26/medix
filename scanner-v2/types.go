package scannerV2

type ScanOptions struct {
	Mode         string   `json:"mode"`        // "dirs" or "files"
	Depth        int      `json:"depth"`       // -1 = unlimited, 0 = root only
	Exts         []string `json:"ext"`         // filters like .mkv, .mp4
	Exclude      []string `json:"exclude"`     // exclude path prefixes
	ShowProgress bool     `json:"progress"`    // show progress bar
	Concurrency  int      `json:"concurrency"` // worker count
	OnlyLeaf     bool     `json:"onlyLeaf"`    // new: only include leaf dirs
	LeafDepth    int      `json:"leafDepth"`   // NEW: 0 = default, 1 = leaf-1, 2 = leaf-2, etc.
	SkipEmpty    bool     `json:"skipEmpty"`   // new: skip empty directories entirely
	Verbose      bool     `json:"verbose"`     // log visited/skipped folders
}

type ScanEntry struct {
	ItemPath   string   `json:"itemPath"`             // Required
	ItemName   string   `json:"itemName"`             // Required
	GroupLabel []string `json:"groupLabel,omitempty"` // Optional
	GroupPath  string   `json:"groupPath,omitempty"`  // Optional
	ItemSize   *int64   `json:"itemSize,omitempty"`   // Optional
	SubEntries []string `json:"subEntries,omitempty"` // Optional
}

type ScanOutput struct {
	GeneratedAt   string      `json:"generated_at"`   // ISO8601 timestamp
	SourcePath    string      `json:"source_path"`    // Cleaned absolute input path
	Mode          string      `json:"mode"`           // "files" or "dirs"
	ItemCount     int         `json:"item_count"`     // len(Items)
	ExcludedCount int         `json:"excluded_count"` // for verbosity/debug
	Duration      string      `json:"duration"`       // Elapsed time
	Items         []ScanEntry `json:"items"`
}
