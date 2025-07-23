// internal/scan/config.go
package scanner

import (
	"fmt"
	"os"
	"sync"
)

type ScanFileConfig struct {
	Concurrency int          `json:"concurrency,omitempty"` // 👈 add this
	Scan        []ScanConfig `json:"scan"`
}

type ScanConfig struct {
	Name        string        `json:"name"`              // e.g. "movies.todo"
	Type        string        `json:"type"`              // e.g. "media", "icon"
	ContentType string        `json:"contentType"`       // e.g. "movies", "tv"
	Phase       string        `json:"phase"`             // e.g. "raw", "staged", "organized"
	Include     []ScanInclude `json:"include"`           // (optional) for future use
	Exclude     []string      `json:"exclude"`           // (optional) for future use
	Output      string        `json:"output"`            // Output file path
	Options     *ScanOptions  `json:"options,omitempty"` // optional overrides
}

type ScanInclude struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

type ScanStrategy interface {
	Scan(sources map[string]string, opts ScanOptions) (any, error) // returns model.MovieOutput or model.TVShowOutput
}

type dirCache struct {
	mu    sync.Mutex
	cache map[string][]os.DirEntry
	// m sync.Map // map[string][]os.DirEntry
}

func (c *dirCache) GetOrRead(path string) ([]os.DirEntry, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entries, ok := c.cache[path]; ok {
		// ✅ Return cached result
		return entries, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", path, err)
		return nil, err
	}

	c.cache[path] = entries
	return entries, nil
}

type ScanMode string

const (
	ScanDirs  ScanMode = "dirs"
	ScanFiles ScanMode = "files"
)

func (m ScanMode) String() string {
	return string(m)
}

type ScanPhase string

const (
	PhaseRaw    ScanPhase = "raw"
	PhaseStaged ScanPhase = "staged"
	PhaseFinal  ScanPhase = "final"
)

func ImpliedScanMode(contentType string, phase ScanPhase) ScanMode {
	switch contentType {
	case "media":
		switch phase {
		case PhaseRaw:
			return ScanFiles
		case PhaseStaged, PhaseFinal:
			return ScanDirs
		default:
			return ScanDirs
		}
	case "icon":
		return ScanFiles // always files for icons
	default:
		return ScanDirs
	}
}

type ScanOptions struct {
	Mode         ScanMode `json:"mode"`                   // "files" or "dirs"
	Depth        int      `json:"depth"`                  // e.g. 4 for raw
	Exts         []string `json:"exts,omitempty"`         // file extensions to include (if files)
	Concurrency  int      `json:"concurrency,omitempty"`  // override global
	ShowProgress bool     `json:"showProgress,omitempty"` // show scan progress
}

type ScanEntry struct {
	Source     string
	GroupLabel []string
	GroupPath  string
	ItemPath   string
	ItemName   string
	ItemSize   *int64
	SubEntries []os.DirEntry
}
