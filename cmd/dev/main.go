package main

import (
	"log"

	"github.com/mrizkifadil26/medix/server"
	"github.com/mrizkifadil26/medix/webgen"
)

func main() {
	log.Println("🛠 Dev mode enabled")

	if err := webgen.GenerateSite("data", "dist"); err != nil {
		log.Fatalf("❌ Initial site generation failed: %v", err)
	}

	go server.WatchAndBuild()
	go server.OpenBrowser("http://localhost:8080")

	if err := server.Serve("dist", "8080"); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
