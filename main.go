package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type MediaMap map[string][]string
type TitleInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func loadJSON(path string) (map[string][]TitleInfo, error) {
	data := make(map[string][]TitleInfo)
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(f, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	loadEnv()
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Println("⚠️  BASE_URL is empty. Assuming local build.")
	}

	router := gin.Default()
	router.Static("/static", "./public")
	// r.Static("/data", "./output")
	router.HTMLRender = loadTemplates()

	router.GET("/", func(c *gin.Context) {
		progress, err := loadProgress("data/progress.json")
		if err != nil {
			c.String(500, "Failed to load progress: %v", err)
			return
		}

		// generate bar string
		barWidth := 50
		filled := (progress.Done * barWidth) / progress.Total
		bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)

		c.HTML(http.StatusOK, "index", gin.H{
			"Genres":  progress.Genres,
			"Done":    progress.Done,
			"Total":   progress.Total,
			"Percent": progress.Percent,
			"Bar":     bar,
		})
	})

	router.GET("/movies", func(c *gin.Context) {
		data, _ := loadJSON("./data/movies.json")
		var genres []string
		for genre := range data {
			genres = append(genres, genre)
		}

		sort.Strings(genres)

		c.HTML(http.StatusOK, "genres", gin.H{
			"BaseURL": baseURL,
			"Type":    "Movies",
			"TypeURL": "movies",
			"Genres":  genres,
		})
	})

	router.GET("/tvshows", func(c *gin.Context) {
		data, _ := loadJSON("./data/tv_shows.json")
		var genres []string
		for genre := range data {
			genres = append(genres, genre)
		}

		sort.Strings(genres)

		c.HTML(http.StatusOK, "genres", gin.H{
			"BaseURL": baseURL,
			"Type":    "TV Shows",
			"TypeURL": "tvshows",
			"Genres":  genres,
		})
	})

	router.GET("/movies/:genre", func(c *gin.Context) {
		genre := c.Param("genre")
		data, _ := loadJSON("./data/movies.json")
		titles := data[genre]
		var formatted []map[string]string

		for _, t := range titles {
			formatted = append(formatted, map[string]string{
				"name":   t.Name,
				"status": t.Status,
			})
		}

		c.HTML(http.StatusOK, "titles", gin.H{
			"BaseURL": baseURL,
			"Type":    "Movies",
			"TypeURL": "movies",
			"Genre":   genre,
			"Titles":  formatted,
		})
	})

	router.GET("/tvshows/:genre", func(c *gin.Context) {
		genre := c.Param("genre")
		data, _ := loadJSON("./data/tv_shows.json")
		titles := data[genre]
		var formatted []map[string]string

		for _, t := range titles {
			formatted = append(formatted, map[string]string{
				"name":   t.Name,
				"status": t.Status,
			})
		}

		c.HTML(http.StatusOK, "titles", gin.H{
			"BaseURL": baseURL,
			"Type":    "TV Shows",
			"TypeURL": "tvshows",
			"Genre":   genre,
			"Titles":  formatted,
		})
	})

	router.Run(":8080")
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// Define a page and its layout
	layout := "templates/layouts/base.tmpl"
	pages, _ := filepath.Glob("templates/pages/*.tmpl")

	for _, page := range pages {
		name := filepath.Base(page[:len(page)-len(filepath.Ext(page))]) // "index" from index.tmpl
		log.Printf("Loading template: %s with layout: %s", name, layout)
		r.AddFromFiles(name, layout, page)
	}

	return r
}

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

func loadProgress(path string) (*Progress, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var progress Progress
	dec := json.NewDecoder(f)
	if err := dec.Decode(&progress); err != nil {
		return nil, err
	}
	return &progress, nil
}

func loadEnv() {
	// Try .env.github first (for CI/deploy), fallback to .env (local)
	if err := godotenv.Load(".env.github"); err != nil {
		_ = godotenv.Load(".env")
	}
}
