package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/joho/godotenv"
)

type Title struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type GenreStatus struct {
	Genre  string
	Raw    int
	Png    int
	Ico    int
	Status string
}

type Progress struct {
	Genres  []GenreStatus
	Done    int
	Total   int
	Percent int
}

type GenreMap map[string][]Title

func mustCreate(path string) *os.File {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", path, err)
	}
	return f
}

func loadProgress(path string) Progress {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to read progress JSON: %v", err)
	}
	defer f.Close()
	var progress Progress
	if err := json.NewDecoder(f).Decode(&progress); err != nil {
		log.Fatalf("Failed to decode progress JSON: %v", err)
	}
	return progress
}

func loadGenreTitleMap(path string) GenreMap {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open genre/title JSON: %v", err)
	}
	defer f.Close()
	var data GenreMap
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		log.Fatalf("Failed to decode genre/title JSON: %v", err)
	}
	return data
}

func main() {
	_ = godotenv.Load()

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = ""
	}

	distDir := "dist"
	moviesDir := filepath.Join(distDir, "movies")

	// Create output directories
	if err := os.MkdirAll(moviesDir, 0755); err != nil {
		log.Fatalf("Failed to create dist directories: %v", err)
	}

	// Parse templates
	tmplIndex := template.Must(template.ParseFiles(
		"templates/layouts/base.tmpl",
		"templates/pages/index.tmpl",
	))
	tmplGenres := template.Must(template.ParseFiles(
		"templates/layouts/base.tmpl",
		"templates/pages/genres.tmpl",
	))
	tmplTitles := template.Must(template.ParseFiles(
		"templates/layouts/base.tmpl",
		"templates/pages/titles.tmpl",
	))

	// Load and generate progress index.html
	progress := loadProgress("data/progress.json")
	fIndex := mustCreate(filepath.Join(distDir, "index.html"))
	defer fIndex.Close()

	err := tmplIndex.ExecuteTemplate(fIndex, "base.tmpl", map[string]any{
		"BaseURL": baseURL,
		"Genres":  progress.Genres,
		"Done":    progress.Done,
		"Total":   progress.Total,
		"Percent": progress.Percent,
	})
	if err != nil {
		log.Fatalf("Failed to render index.html: %v", err)
	}

	// Load genre-title map for movies
	genreMap := loadGenreTitleMap("data/movies.json")

	// Extract sorted genres list
	var genres []string
	for g := range genreMap {
		genres = append(genres, g)
	}
	sort.Strings(genres)

	// Generate movies genre index page
	fGenres := mustCreate(filepath.Join(moviesDir, "index.html"))
	defer fGenres.Close()

	err = tmplGenres.ExecuteTemplate(fGenres, "base.tmpl", map[string]any{
		"BaseURL": baseURL,
		"Type":    "Movies",
		"TypeURL": "movies",
		"Genres":  genres,
	})
	if err != nil {
		log.Fatalf("Failed to render movies genre index: %v", err)
	}

	// Generate title pages per genre
	for _, genre := range genres {
		titles := genreMap[genre]
		outPath := filepath.Join(moviesDir, genre+".html")

		f, err := os.Create(outPath)
		if err != nil {
			log.Printf("Failed to create %s: %v", outPath, err)
			continue
		}

		var formatted []map[string]string
		for _, t := range titles {
			formatted = append(formatted, map[string]string{
				"name":   t.Name,
				"status": t.Status,
			})
		}

		err = tmplTitles.ExecuteTemplate(f, "base.tmpl", map[string]any{
			"BaseURL": baseURL,
			"Genre":   genre,
			"Titles":  formatted,
			"Type":    "Movies",
			"TypeURL": "movies",
		})
		f.Close()

		if err != nil {
			log.Printf("Failed to render titles page for %s: %v", genre, err)
		}
	}

	log.Println("✅ Static site generated in ./dist")
}
