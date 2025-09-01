package server

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mrizkifadil26/medix/utils/logger"
	"github.com/mrizkifadil26/medix/webgen"
)

func WatchAndBuild() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("Watcher init failed: " + err.Error())
		return
	}
	defer watcher.Close()

	paths := []string{"data", "templates", "assets"}
	for _, path := range paths {
		_ = watcher.Add(path)
		logger.Info("👁️ Watching: " + path)
	}

	debounce := time.Now()

	for {
		select {
		case event := <-watcher.Events:
			if time.Since(debounce) < 300*time.Millisecond {
				continue
			}
			debounce = time.Now()

			logger.Info("🔁 Change detected: " + event.Name)
			err := webgen.GenerateSite("data", "dist")
			if err != nil {
				// log.Printf("❌ Generate error: %v", err)
				logger.Error("❌ Generate error: " + err.Error())

			} else {
				// log.Println("✅ Site regenerated.")
				logger.Info("✅ Site regenerated.")
			}

		case err := <-watcher.Errors:
			// log.Println("Watcher error:", err)
			logger.Error("❌ Watcher error: " + err.Error())
		}
	}
}
