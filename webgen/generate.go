package webgen

import (
	"path/filepath"

	"github.com/mrizkifadil26/medix/logger"
)

func GenerateSite(inputDir, outputDir string) error {
	// Paths
	dataPath := inputDir
	outData := filepath.Join(outputDir, "data")

	// Step 1: Render static pages
	logger.Step("🎨 Rendering static pages...")
	RenderStaticPages()

	// Step 2: Render data-driven pages
	logger.Step("📦 Rendering data pages...")
	RenderDataPage("movies.json", "Movies", "movies.html")
	RenderDataPage("tv_shows.json", "TV Shows", "tv.html")
	// Add more data pages if needed

	// Step 3: Copy raw JSONs to output dir
	logger.Step("📁 Copying data to dist/data...")
	if err := CopyDir(dataPath, outData); err != nil {
		return err
	}

	// Step 4: Copy static assets
	logger.Step("🎨 Copying assets to dist/assets...")
	if err := CopyDir("assets", outputDir); err != nil {
		return err
	}

	logger.Done("✅ Site generation complete.")
	return nil
}
