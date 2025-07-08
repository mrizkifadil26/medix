package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/mrizkifadil26/medix/model"
)

const maxConcurrency = 8

type dirCache struct {
	m sync.Map // map[string][]os.DirEntry
}

func ScanDirectory(contentType, root string) model.RawOutput {
	fmt.Printf("Starting scan for type: %s in root: %s\n", contentType, root)
	result := model.RawOutput{
		Type:        contentType,
		GeneratedAt: time.Now(),
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Failed to read root directory: %v\n", err)
		return result
	}

	var genres []os.DirEntry
	for _, g := range entries {
		if g.IsDir() {
			genres = append(genres, g)
		}
	}

	fmt.Printf("Found %d genre directories\n", len(genres))

	bar := progressbar.NewOptions(len(genres),
		progressbar.OptionSetDescription("Scanning genres"),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
	)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
	genreMap := make(map[string]model.GenreBlock)

	cache := &dirCache{}

	for _, genre := range genres {
		wg.Add(1)
		sem <- struct{}{}

		go func(genre os.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()
			defer bar.Add(1)

			genreName := genre.Name()
			genrePath := filepath.Join(root, genreName)

			items := scanGenre(genrePath, cache, contentType)
			if len(items) == 0 {
				return
			}

			mu.Lock()
			genreMap[genreName] = model.GenreBlock{
				Genre: genreName,
				Items: items,
			}
			mu.Unlock()
		}(genre)
	}

	wg.Wait()

	// Reconstruct result.Data in sorted order
	for _, genre := range genres {
		if block, ok := genreMap[genre.Name()]; ok {
			result.Data = append(result.Data, block)
		}
	}

	fmt.Println("Scan complete.")

	return result
}

func scanGenre(genrePath string, cache *dirCache, contentType string) []model.RawItem {
	entries := cache.Read(genrePath)
	if entries == nil {
		return nil
	}

	// Sort title folders
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var items []model.RawItem

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		titlePath := filepath.Join(genrePath, entry.Name())
		subEntries := cache.Read(titlePath)
		if subEntries == nil {
			continue
		}

		var children []model.RawChild
		switch contentType {
		case "movies":
			children = extractChildren(titlePath, subEntries, cache)
		case "tvshows":
			children = listSeasons(titlePath, subEntries)
		}

		itemType := "single"
		if contentType == "movies" && len(children) > 0 {
			itemType = "collection"
		}

		status := resolveStatus(subEntries)
		items = append(items, model.RawItem{
			Type:     itemType,
			Name:     entry.Name(),
			Path:     titlePath,
			Status:   status,
			Children: children,
		})
	}

	return items
}

func extractChildren(parent string, entries []os.DirEntry, cache *dirCache) []model.RawChild {
	// Sort child directories
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var children []model.RawChild
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		childPath := filepath.Join(parent, e.Name())
		subEntries := cache.Read(childPath)
		if subEntries == nil {
			continue
		}

		status := resolveStatus(subEntries)
		children = append(children, model.RawChild{
			Name:   e.Name(),
			Path:   childPath,
			Status: status,
		})
	}

	return children
}

func listSeasons(parent string, entries []os.DirEntry) []model.RawChild {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var children []model.RawChild
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		childPath := filepath.Join(parent, e.Name())
		children = append(children, model.RawChild{
			Name:   e.Name(),
			Path:   childPath,
			Status: "unknown",
		})
	}

	return children
}

func resolveStatus(entries []os.DirEntry) string {
	hasIco := false
	hasIni := false

	for _, f := range entries {
		if f.IsDir() {
			continue
		}
		switch filepath.Ext(f.Name()) {
		case ".ico":
			hasIco = true
		case ".ini":
			if f.Name() == "desktop.ini" {
				hasIni = true
			}
		}

		// Early exit once both are found
		if hasIco && hasIni {
			return "ok"
		}
	}

	if hasIco {
		return "warn"
	}

	return "missing"
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
