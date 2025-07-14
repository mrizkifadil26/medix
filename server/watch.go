package server

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mrizkifadil26/medix/webgen"
)

func WatchAndBuild() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	paths := []string{"data", "templates", "assets"}
	for _, path := range paths {
		_ = watcher.Add(path)
	}

	debounce := time.Now()

	for {
		select {
		case event := <-watcher.Events:
			if time.Since(debounce) < 300*time.Millisecond {
				continue
			}
			debounce = time.Now()

			log.Printf("🔁 Change detected: %s", event.Name)
			err := webgen.GenerateSite("data", "dist")
			if err != nil {
				log.Printf("❌ Generate error: %v", err)
			} else {
				log.Println("✅ Site regenerated.")
			}

		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}
