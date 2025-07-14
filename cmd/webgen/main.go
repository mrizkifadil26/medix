package main

import (
	"flag"
	"log"

	"github.com/mrizkifadil26/medix/webgen"
)

var distDirs = []string{
	"dist/css",
	"dist/js",
	"dist/data",
}

func main() {
	// Optional flags for future extensibility
	clean := flag.Bool("clean", true, "Clean the dist directory before build")
	verbose := flag.Bool("v", false, "Enable verbose output")

	flag.Parse()

	if *verbose {
		log.Println("🔨 Building static site with templates...")
	}

	if *clean {
		webgen.CleanDist()
	}

	webgen.EnsureDirs(distDirs...)
	webgen.CopyAssets()
	webgen.RenderPages()

	if *verbose {
		log.Println("✅ Static site generated in dist/")
	}
}
