// internal/scan/config.go
package scan

import (
	"fmt"
	"os"
	"sync"
)

type ScanConfig struct {
	ContentType string   `json:"content_type"`       // "movies" or "tvshows"
	Sources     []string `json:"sources"`            // List of directories
	OutputPath  string   `json:"output_path"`        // Output file path
	Strategy    string   `json:"strategy,omitempty"` // (optional) for future use
}

type ScanConfigFile struct {
	Concurrency int          `json:"concurrency,omitempty"` // 👈 add this
	Configs     []ScanConfig `json:"configs"`
}

type ScanStrategy[T any] interface {
	Scan(roots []string) T // returns model.MovieOutput or model.TVShowOutput
}

type dirCache struct {
	m sync.Map // map[string][]os.DirEntry
}

func (dc *dirCache) Read(path string) []os.DirEntry {
	if val, ok := dc.m.Load(path); ok {
		return val.([]os.DirEntry)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", path, err)
		return nil
	}

	dc.m.Store(path, entries)
	return entries
}
