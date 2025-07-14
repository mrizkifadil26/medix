package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/webgen"
)

func main() {
	var inputDir, outputDir string
	flag.StringVar(&inputDir, "input", "data", "Input data directory")
	flag.StringVar(&outputDir, "output", "dist", "Output directory")
	flag.Parse()

	log.Println("⚙️ Starting site generation...")
	err := webgen.GenerateSite(inputDir, outputDir)
	if err != nil {
		log.Fatalf("❌ Generation failed: %v", err)
	}

	fmt.Println("✅ Static site successfully generated.")
}
