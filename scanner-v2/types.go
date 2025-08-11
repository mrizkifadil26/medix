package scannerV2

// type ScanOptions = Options

type ScanEntry struct {
	Path    string `json:"path"`     // Absolute path
	RelPath string `json:"rel_path"` // Relative to scan root
	Name    string `json:"name"`     // Basename
	Type    string `json:"type"`     // "file", "dir", "symlink", etc

	Ext     string `json:"ext,omitempty"`      // ".mkv", ".txt", etc
	Size    *int64 `json:"size,omitempty"`     // Optional size
	ModTime string `json:"mod_time,omitempty"` // ISO8601

	GroupLabel    []string    `json:"group_label,omitempty"`    // e.g. ["Action"], ["Action", "Marvel"]
	GroupPath     string      `json:"group_path,omitempty"`     // e.g. "Action/Marvel"
	AncestorPaths []string    `json:"ancestor_paths,omitempty"` // ["Action", "Action/MCU"]
	Children      []ScanEntry `json:"children,omitempty"`       // Recursive entries
}

type ScanOutput struct {
	Version     string        `json:"version"`            // Schema version, e.g. "1.0.0"
	GeneratedAt string        `json:"generated_at"`       // ISO8601 timestamp
	SourcePath  string        `json:"source_path"`        // Absolute root path scanned
	Mode        string        `json:"mode"`               // "files" | "dirs" | "mixed"
	ItemCount   int           `json:"item_count"`         // len(Items)
	DurationMs  int64         `json:"duration_ms"`        // Total elapsed time in milliseconds
	Tags        []string      `json:"tags,omitempty"`     // Optional job/context tags
	Stats       *WalkOptions  `json:"stats,omitempty"`    // Deep stats from walker
	Errors      []ScanError   `json:"errors,omitempty"`   // Errors encountered (path + reason)
	Warnings    []ScanWarning `json:"warnings,omitempty"` // Non-critical issues
	Items       []ScanEntry   `json:"items"`              // Final matched items
}

type ScanError struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type ScanWarning struct {
	Path   string `json:"path"`
	Detail string `json:"detail"`
}

type SubentriesMode string

const (
	SubentriesNone   SubentriesMode = "none"
	SubentriesFlat   SubentriesMode = "flat"
	SubentriesNested SubentriesMode = "nested"
	SubentriesAuto   SubentriesMode = "auto"
)
