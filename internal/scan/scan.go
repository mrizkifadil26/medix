package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

const maxConcurrency = 8

type dirCache struct {
	m sync.Map // map[string][]os.DirEntry
}

func ScanAll(cfg ScanConfig) model.RawOutput {
	result := model.RawOutput{
		Type:        "raw",
		GeneratedAt: time.Now(),
	}

	cache := &dirCache{}
	genreMap := make(map[string]model.RawGenre)
	for _, root := range cfg.Sources {
		scanSingleRoot(cfg.ContentType, root, genreMap, cache)
	}

	// Sort and reassemble
	var genreNames []string
	for name := range genreMap {
		genreNames = append(genreNames, name)
	}
	sort.Strings(genreNames)

	for _, name := range genreNames {
		result.Data = append(result.Data, genreMap[name])
	}

	fmt.Println("All sources scanned.")
	return result
}

func scanSingleRoot(contentType, root string, genreMap map[string]model.RawGenre, cache *dirCache) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Failed to read root directory: %v\n", err)
		return
	}

	var genres []os.DirEntry
	for _, g := range entries {
		if g.IsDir() {
			genres = append(genres, g)
		}
	}

	fmt.Printf("Scanning %d genres in %s\n", len(genres), root)
	bar := progressbar.NewOptions(len(genres),
		progressbar.OptionSetDescription("Scanning genres"),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
	)

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
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
			existing := genreMap[genreName]
			existing.Name = genreName
			existing.Items = append(existing.Items, items...)
			genreMap[genreName] = existing
			mu.Unlock()
		}(genre)
	}

	wg.Wait()
}

func scanGenre(genrePath string, cache *dirCache, contentType string) []model.RawEntry {
	entries := cache.Read(genrePath)
	if entries == nil {
		return nil
	}

	// Sort title folders
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var items []model.RawEntry

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		titlePath := filepath.Join(genrePath, entry.Name())
		subEntries := cache.Read(titlePath)
		if subEntries == nil {
			continue
		}

		status := resolveStatus(subEntries)
		ico := findIcon(titlePath, subEntries)

		list := extractChildren(titlePath, subEntries, cache, contentType)
		itemType := "single"
		if contentType == "movies" && len(list) > 0 {
			itemType = "collection"
		}

		rawEntry := model.RawEntry{
			Type:   itemType,
			Name:   entry.Name(),
			Path:   titlePath,
			Status: status,
			Icon:   ico,
		}

		switch contentType {
		case "movies":
			if len(list) > 0 {
				rawEntry.Items = &model.RawEntryItems{Entries: list}
			}
		case "tvshows":
			seasons := extractSeasonNames(subEntries)
			if len(seasons) > 0 {
				rawEntry.Items = &model.RawEntryItems{Seasons: seasons}
			}
		}

		items = append(items, rawEntry)
	}

	return items
}

func extractChildren(parent string, entries []os.DirEntry, cache *dirCache, contentType string) []model.RawEntry {
	// Sort child directories
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var children []model.RawEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		childPath := filepath.Join(parent, e.Name())
		subEntries := cache.Read(childPath)
		if subEntries == nil {
			continue
		}

		childType := "single"
		for _, sub := range subEntries {
			if sub.IsDir() {
				childType = "collection"
				break
			}
		}

		status := resolveStatus(subEntries)
		ico := findIcon(childPath, subEntries)
		children = append(children, model.RawEntry{
			Type:   childType,
			Name:   e.Name(),
			Path:   childPath,
			Status: status,
			Icon:   ico,
		})
	}

	return children
}

func extractSeasonNames(entries []os.DirEntry) []string {
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
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

func findIcon(dir string, entries []os.DirEntry) *model.IconMeta {
	for _, f := range entries {
		if f.IsDir() || filepath.Ext(f.Name()) != ".ico" {
			continue
		}

		info, err := os.Stat(filepath.Join(dir, f.Name()))
		if err != nil {
			continue
		}

		return &model.IconMeta{
			ID:       util.Slugify(f.Name()), // Use the file name as ID
			Name:     f.Name(),
			FullPath: filepath.Join(dir, f.Name()),
			Size:     info.Size(),
		}
	}
	return nil
}
