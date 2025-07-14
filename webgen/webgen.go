package webgen

import (
	"log"
	"os"
	"path/filepath"
)

func CleanDist() {
	err := os.RemoveAll("dist")
	if err != nil {
		log.Fatal(err)
	}
}

func CopyAssets() {
	err := CopyDir("assets/js", "dist/js")
	Must(err)

	err = CopyDir("assets/css", "dist/css")
	Must(err)

	// Copy root-level icons
	files := []string{
		"assets/favicon.ico",
		"assets/apple-touch-icon.png",
		"assets/android-chrome-192x192.png",
		"assets/android-chrome-512x512.png",
	}
	for _, file := range files {
		err := CopyFile(file, "dist/"+filepath.Base(file))
		Must(err)
	}

	err = CopyDir("data", "dist/data")
	Must(err)
}

func RenderPages() {
	RenderStaticPages()
	RenderDataPage("movies.json", "Movies", "movies.html")
	RenderDataPage("tv_shows.json", "TV Shows", "tvshows.html")
}
