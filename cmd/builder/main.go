package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func renderPageFiles(files []string, outPath string, data any) {
	tmpl, err := template.ParseFiles(files...)
	must(err)

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "base", data)
	must(err)

	// 🔍 Debug: print to console
	// fmt.Println("----- " + outPath + " -----")
	// fmt.Println(buf.String())

	// 💾 Write to output file
	f, err := os.Create(outPath)
	must(err)
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	must(err)
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
		targetPath := filepath.Join(dst, path[len(src):])
		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}
		return copyFile(path, targetPath)
	})
}

func main() {
	log.Println("🔨 Building static site with templates...")

	// Clean and prepare dist
	os.RemoveAll("dist")
	os.MkdirAll("dist/css", 0755)
	os.MkdirAll("dist/js", 0755)
	os.MkdirAll("dist/data", 0755)

	// Copy static assets
	copyDir("public/js", "dist/js")
	copyDir("public/css", "dist/css")
	copyDir("data", "dist/data")

	// Render Home
	renderPageFiles([]string{
		"templates/layouts/base.html",
		"templates/pages/index.html",
	}, "dist/index.html", map[string]any{
		"Title": "Media Tracker",
	})

	// Render About
	renderPageFiles([]string{
		"templates/layouts/base.html",
		"templates/pages/about.html",
	}, "dist/about.html", map[string]any{
		"Title": "About",
	})

	// Render Movies and TV Shows
	renderData := func(jsonFile, title, outFile string) {
		raw, err := os.ReadFile("data/" + jsonFile)
		must(err)

		var obj any
		must(json.Unmarshal(raw, &obj))

		pageData := map[string]any{
			"Title": title,
			"Data":  template.JS(string(raw)), // JS object injection
		}

		renderPageFiles([]string{
			"templates/layouts/base.html",
			"templates/pages/category.html",
		}, "dist/"+outFile, pageData)
	}

	renderData("movies.json", "Movies", "movies.html")
	renderData("tv_shows.json", "TV Shows", "tvshows.html")

	log.Println("✅ Static site generated in dist/")
}
