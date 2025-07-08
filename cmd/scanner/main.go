package main

import (
	"net/http"
	_ "net/http/pprof" // Enable pprof
	"time"

	"fmt"
	"log"
	"os"

	"github.com/mrizkifadil26/medix/internal/scanner"
	"github.com/mrizkifadil26/medix/util"
)

func main() {
	go func() {
		log.Println("📊 pprof running at http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))

		fmt.Println("⏳ Waiting 30 seconds for pprof inspection...")
		time.Sleep(30 * time.Second)
	}()

	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run scanner.go <movies|tvshows>")
	}

	mode := os.Args[1]
	var rootDir, outputPath string

	switch mode {
	case "movies":
		rootDir = "/mnt/e/Media/Movies"
		outputPath = "data/movies.raw.json"
	case "tvshows":
		rootDir = "/mnt/e/Media/TV Shows"
		outputPath = "data/tv_shows.raw.json"
	default:
		log.Fatalf("Unknown mode: %s", mode)
	}

	result := scanner.ScanDirectory(mode, rootDir)
	if len(result.Data) == 0 {
		log.Printf("⚠️ No entries found in %s\n", rootDir)
		return
	}

	if err := util.WriteJSON(outputPath, result); err != nil {
		log.Fatalf("❌ Failed to write JSON: %v", err)
	}
	fmt.Printf("✅ %s written (%d genres)\n", outputPath, len(result.Data))
}
