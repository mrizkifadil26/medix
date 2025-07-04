package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type GenreStatus struct {
	Genre  string `json:"genre"`
	Raw    int    `json:"raw"`
	Png    int    `json:"png"`
	Ico    int    `json:"ico"`
	Status string `json:"status"`
}

type Progress struct {
	Genres  []GenreStatus `json:"genres"`
	Done    int           `json:"done"`
	Total   int           `json:"total"`
	Percent int           `json:"percent"`
}

func main() {
	baseDir := "/mnt/c/Users/Rizki/OneDrive/Pictures/Icons/Personal Icon Pack/Movies"
	outputPath := flag.String("out", "data/progress.json", "Output JSON file")
	flag.Parse()

	formats := []string{"RAW", "PNG", "ICO"}
	formatDirs := map[string]string{}
	for _, format := range formats {
		formatDirs[format] = filepath.Join(baseDir, format)
	}

	fmt.Printf("Scanning directories: %s\n", formatDirs)

	genreSet := make(map[string]bool)
	table := make(map[string]map[string]int)

	for _, format := range formats {
		base := formatDirs[format]
		if _, err := os.Stat(base); os.IsNotExist(err) {
			continue
		}

		entries, _ := os.ReadDir(base)
		for _, e := range entries {
			if !e.IsDir() || e.Name() == "Collection" {
				continue
			}
			genre := e.Name()
			genreSet[genre] = true
			genrePath := filepath.Join(base, genre)

			count := 0
			filepath.WalkDir(genrePath, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() && strings.Contains(path, "Collection") {
					return filepath.SkipDir
				}
				if !d.IsDir() {
					count++
				}
				return nil
			})
			if _, ok := table[genre]; !ok {
				table[genre] = make(map[string]int)
			}
			table[genre][format] = count
		}
	}

	var genres []string
	for g := range genreSet {
		genres = append(genres, g)
	}
	sort.Slice(genres, func(i, j int) bool {
		return table[genres[i]]["RAW"] > table[genres[j]]["RAW"]
	})

	fmt.Printf("\033[1;32m%-15s%8s%8s%8s%10s\033[0m\n", "Genre", "RAW", "PNG", "ICO", "Completed")

	var progress Progress
	for _, genre := range genres {
		raw := table[genre]["RAW"]
		png := table[genre]["PNG"]
		ico := table[genre]["ICO"]
		status := "✅"
		if raw != png || raw != ico {
			status = "⏳"
		} else {
			progress.Done++
		}

		fmt.Printf("%-15s%8d%8d%8d%10s\n", genre, raw, png, ico, status)

		progress.Genres = append(progress.Genres, GenreStatus{
			Genre:  genre,
			Raw:    raw,
			Png:    png,
			Ico:    ico,
			Status: status,
		})
	}

	progress.Total = len(progress.Genres)
	if progress.Total > 0 {
		progress.Percent = (progress.Done * 100) / progress.Total
	}
	barWidth := 50
	filled := (progress.Done * barWidth) / progress.Total
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)
	fmt.Printf("\nProgress: [%s] %d%% (%d/%d genres completed)\n",
		bar, progress.Percent, progress.Done, progress.Total)

	os.MkdirAll(filepath.Dir(*outputPath), 0755)
	f, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("❌ Error writing JSON: %v\n", err)
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(progress); err != nil {
		fmt.Printf("❌ Failed to encode JSON: %v\n", err)
		return
	}
	fmt.Printf("✅ JSON written to %s\n", *outputPath)
}
