// internal/scan/config.go
package scanner

import (
	"fmt"
	"os"
	"sync"

	"github.com/mrizkifadil26/medix/model"
)

type ScanFileConfig struct {
	Concurrency int          `json:"concurrency,omitempty"` // 👈 add this
	Scan        []ScanConfig `json:"scan"`
}

type ScanConfig struct {
	Name    string         `json:"name"`    // "movies" or "tv"
	Type    string         `json:"type"`    // List of directories
	Output  string         `json:"output"`  // Output file path
	Include []IncludeEntry `json:"include"` // (optional) for future use
	Exclude []string       `json:"exclude"` // (optional) for future use
}

type IncludeEntry struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

type ScanStrategy interface {
	Scan(sources map[string]string) (model.MediaOutput, error) // returns model.MovieOutput or model.TVShowOutput
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
