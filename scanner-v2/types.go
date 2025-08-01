package scannerV2

type ScanOptions = Options

type ScanEntry struct {
	ItemPath   string      `json:"itemPath"`             // Required
	ItemName   string      `json:"itemName"`             // Required
	GroupLabel []string    `json:"groupLabel,omitempty"` // Optional
	GroupPath  string      `json:"groupPath,omitempty"`  // Optional
	ItemSize   *int64      `json:"itemSize,omitempty"`   // Optional
	SubPaths   []string    `json:"subPaths,omitempty"`   // for "path" mode
	SubEntries []ScanEntry `json:"subEntries,omitempty"` // for "entry" and "recursive" mode
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

type SubentriesMode string

const (
	SubentriesNone   SubentriesMode = "none"
	SubentriesFlat   SubentriesMode = "flat"
	SubentriesNested SubentriesMode = "nested"
	SubentriesAuto   SubentriesMode = "auto"
)
