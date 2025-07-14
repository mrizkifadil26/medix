package webgen

import (
	"log"
	"path/filepath"
)

func GenerateSite(inputDir, outputDir string) error {
	// Paths
	dataPath := inputDir
	outData := filepath.Join(outputDir, "data")
	outAssets := filepath.Join(outputDir, "assets")

	// Step 1: Render static pages
	log.Println("🎨 Rendering static pages...")
	RenderStaticPages()

	// Step 2: Render data-driven pages
	log.Println("📦 Rendering data pages...")
	RenderDataPage("movies.json", "Movies", "movies.html")
	RenderDataPage("tv_shows.json", "TV Shows", "tv.html")
	// Add more data pages if needed

	// Step 3: Copy raw JSONs to output dir
	log.Println("📁 Copying data to dist/data...")
	if err := CopyDir(dataPath, outData); err != nil {
		return err
	}

	// Step 4: Copy static assets
	log.Println("🎨 Copying assets to dist/assets...")
	if err := CopyDir("assets", outAssets); err != nil {
		return err
	}

	log.Println("✅ Site generation complete.")
	return nil
}
