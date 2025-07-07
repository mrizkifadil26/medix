package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
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
	fmt.Println("----- " + outPath + " -----")
	fmt.Println(buf.String())

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
	watch := flag.Bool("watch", false, "Watch files and rebuild on changes")
	flag.Parse()

	if *watch {
		startWatcher()
	} else {
		buildSite()
	}
}

func buildSite() {
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

func startWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watchDirs := []string{"templates", "data", "public"}
	for _, dir := range watchDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("👀 Watching for changes...")

	var timer *time.Timer

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				log.Printf("🔁 Change detected: %s", event.Name)

				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(500*time.Millisecond, func() {
					buildSite()
				})
			}

		case err := <-watcher.Errors:
			log.Println("❌ Watcher error:", err)
		}
	}
}
