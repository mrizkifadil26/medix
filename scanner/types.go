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

type ScanOptions struct {
	Mode         ScanMode
	Depth        int
	Exts         []string // e.g. []string{".mkv", ".mp4"}
	Concurrency  int
	ShowProgress bool
}

type ScannedItem struct {
	GroupLabel string
	GroupPath  string
	ItemPath   string
	ItemName   string
	SubEntries []os.DirEntry
}
