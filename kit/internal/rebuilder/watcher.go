package rebuilder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
)

var watchExtensions string

func init() {
	pflag.StringVar(&watchExtensions, "watch.extensions", ".go,.css,.js", "comma separated extensions to watch for recompile")
}

type Watcher interface {
	Watch(reload chan bool) error
}

func watcher() Watcher {
	return &fileWatcher{
		watcher:    nil,
		extensions: strings.Split(watchExtensions, ","),
	}
}

type fileWatcher struct {
	watcher    *fsnotify.Watcher
	extensions []string
}

func (f *fileWatcher) Watch(reload chan bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		return watcher.AddWith(path)
	})

	if err != nil {
		return fmt.Errorf("error loading paths: %w", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if !slices.Contains(f.extensions, filepath.Ext(event.Name)) {
					continue
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) ||
					event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					reload <- true
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	<-make(chan struct{})

	return nil
}
