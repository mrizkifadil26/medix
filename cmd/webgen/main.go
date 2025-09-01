package main

import (
	"flag"

	"github.com/mrizkifadil26/medix/utils/logger"
	"github.com/mrizkifadil26/medix/webgen"
)

func main() {
	var (
		inputDir  string
		outputDir string
		dryRun    bool
	)

	flag.StringVar(&inputDir, "input", "data", "Input directory")
	flag.StringVar(&outputDir, "output", "dist", "Output directory")
	flag.BoolVar(&dryRun, "dry", false, "Dry-run mode (no output written)")
	flag.Parse()

	webgen.DryRun = dryRun

	logger.Info("⚙️ Starting site generation...")
	err := webgen.GenerateSite(inputDir, outputDir)
	if err != nil {
		logger.Error("❌ Generation failed: " + err.Error())
	}
}
