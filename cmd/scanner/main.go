package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type TitleInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "ok", "warn", or "missing"
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

	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("⚠️ Failed to read: %s\n", root)
		return result
	}

	for _, genreDir := range entries {
		if !genreDir.IsDir() {
			continue
		}

		genrePath := filepath.Join(root, genreDir.Name())
		titleDirs, _ := os.ReadDir(genrePath)

		var titles []TitleInfo
		for _, titleDir := range titleDirs {
			if !titleDir.IsDir() {
				continue
			}

			titlePath := filepath.Join(genrePath, titleDir.Name())

			icoExists := hasIcoFile(titlePath)
			iniExists := fileExists(filepath.Join(titlePath, "desktop.ini"))

			status := "missing"
			if icoExists && iniExists {
				status = "ok"
			} else if icoExists {
				status = "warn"
			}

			titles = append(titles, TitleInfo{
				Name:   titleDir.Name(),
				Status: status,
			})
		}

		if len(titles) > 0 {
			result[genreDir.Name()] = titles
		}
	}

	return result
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
