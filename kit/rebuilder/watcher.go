package rebuilder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// runs the file watcher and notifies the manager when a change
// is detected through the changed channel.
func runWatcher(changed chan bool) {
	// Create new watcher.
	watcher, err := buildWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {
					changed <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Block main goroutine forever.
	<-make(chan struct{})
}

func buildWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, fmt.Errorf("error creating watcher: %w", err)
	}

	// Add all files that need to be watched to the
	// watcher so it notifies the errors that it needs to
	// restart.
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if config.isExcludedPath(path) {
			return nil
		}

		if !config.isWatchedExtension(filepath.Ext(path)) {
			return nil
		}

		return watcher.Add(path)
	})

	if err != nil {
		return watcher, fmt.Errorf("error loading paths: %w", err)
	}

	return watcher, err
}
