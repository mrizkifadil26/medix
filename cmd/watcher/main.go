package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

func runDeployScript() {
	cmd := exec.Command("bash", "./scripts/push-data.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Deploy error: %v\nOutput: %s", err, output)
		return
	}
	log.Println("Deploy successful.")
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	var debounce <-chan time.Time

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
					log.Println("Detected change:", event.Name)
					debounce = time.After(2 * time.Second)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)

			case <-debounce:
				log.Println("Running deploy after change...")
				runDeployScript()
			}
		}
	}()

	err = watcher.Add("data")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching 'data/' for changes...")
	<-done
}
