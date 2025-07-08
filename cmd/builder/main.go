package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ensureDirs(paths ...string) {
	for _, path := range paths {
		must(os.MkdirAll(path, 0755))
	}
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, src)
		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

func renderTemplate(files []string, outPath string, data any) {
	tmpl, err := template.ParseFiles(files...)
	must(err)

	var buf bytes.Buffer
	must(tmpl.ExecuteTemplate(&buf, "base", data))

	f, err := os.Create(outPath)
	must(err)
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	must(err)
}

func renderStaticPages() {
	pages := []struct {
		Files   []string
		OutPath string
		Data    map[string]any
	}{
		{
			Files:   []string{"templates/layouts/base.html", "templates/pages/index.html"},
			OutPath: "dist/index.html",
			Data:    nil,
		},
		{
			Files:   []string{"templates/layouts/base.html", "templates/pages/about.html"},
			OutPath: "dist/about.html",
			Data:    nil,
		},
	}

	for _, page := range pages {
		renderTemplate(page.Files, page.OutPath, page.Data)
	}
}

func renderDataPage(jsonFile, title, outFile string) {
	raw, err := os.ReadFile("data/" + jsonFile)
	must(err)

	var parsed any
	must(json.Unmarshal(raw, &parsed))

	data := map[string]any{
		"Title": title,
		"Type":  strings.TrimSuffix(jsonFile, ".json"),
		"Data":  template.JS(string(raw)),
	}

	renderTemplate(
		[]string{"templates/layouts/base.html", "templates/pages/list.html"},
		"dist/"+outFile,
		data,
	)
}

func main() {
	log.Println("🔨 Building static site with templates...")

	// Clean and prepare dist
	must(os.RemoveAll("dist"))
	ensureDirs("dist/css", "dist/js", "dist/data")

	// Copy static assets
	must(copyDir("public/js", "dist/js"))
	must(copyDir("public/css", "dist/css"))
	must(copyDir("data", "dist/data"))

	// Render Home
	// Render static and data-driven pages
	renderStaticPages()
	renderDataPage("movies.json", "Movies", "movies.html")
	renderDataPage("tv_shows.json", "TV Shows", "tvshows.html")

	log.Println("✅ Static site generated in dist/")
}
