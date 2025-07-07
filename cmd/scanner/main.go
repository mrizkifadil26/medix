package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type TitleInfo struct {
	Name   string      `json:"name"`
	Status string      `json:"status"`
	Group  []TitleInfo `json:"group,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run scanner.go <movies|tvshows>")
	}

	mode := os.Args[1]
	switch mode {
	case "movies":
		scanAndSave("Movies", "/mnt/e/Media/Movies", "data/movies.json")
	case "tvshows":
		scanAndSave("TV Shows", "/mnt/e/Media/TV Shows", "data/tv_shows.json")
	default:
		log.Fatalf("Unknown mode: %s", mode)
	}
}

func scanAndSave(label, dir, out string) {
	data := scanSection(dir)
	if len(data) == 0 {
		log.Printf("⚠️ No entries found in %s\n", dir)
		return
	}
	if err := saveJSON(data, out); err != nil {
		log.Fatalf("❌ Failed to write %s: %v", out, err)
	}
	fmt.Printf("✅ %s written (%d genres)\n", out, len(data))
}

func scanSection(root string) map[string][]TitleInfo {
	result := make(map[string][]TitleInfo)

	genres, err := os.ReadDir(root)
	if err != nil {
		log.Printf("⚠️ Failed to read: %s\n", root)
		return result
	}

	for _, genreDir := range genres {
		if !genreDir.IsDir() {
			continue
		}

		genreName := genreDir.Name()
		genrePath := filepath.Join(root, genreName)

		entries, _ := os.ReadDir(genrePath)
		var titles []TitleInfo

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			entryPath := filepath.Join(genrePath, entry.Name())
			subEntries, _ := os.ReadDir(entryPath)

			// Check if it's a collection (has subfolders)
			hasSubfolders := false
			for _, sub := range subEntries {
				if sub.IsDir() {
					hasSubfolders = true
					break
				}
			}

			if hasSubfolders {
				// Collection: group of titles
				var grouped []TitleInfo
				for _, sub := range subEntries {
					if !sub.IsDir() {
						continue
					}
					titlePath := filepath.Join(entryPath, sub.Name())
					grouped = append(grouped, TitleInfo{
						Name:   sub.Name(),
						Status: resolveStatus(titlePath),
					})
				}
				titles = append(titles, TitleInfo{
					Name:   entry.Name(), // e.g. "John Wick Collection"
					Status: "ok",         // or calculate from group if needed
					Group:  grouped,
				})
			} else {
				// Single title folder
				titles = append(titles, TitleInfo{
					Name:   entry.Name(),
					Status: resolveStatus(entryPath),
				})
			}
		}

		if len(titles) > 0 {
			result[genreName] = titles
		}
	}

	return result
}

func resolveStatus(path string) string {
	ico := hasIcoFile(path)
	ini := fileExists(filepath.Join(path, "desktop.ini"))

	if ico && ini {
		return "ok"
	} else if ico {
		return "warn"
	}
	return "missing"
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func hasIcoFile(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".ico" {
			return true
		}
	}
	return false
}

func saveJSON(data any, path string) error {
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
