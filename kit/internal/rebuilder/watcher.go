package rebuilder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
)

var watchExtensions string

func init() {
	pflag.StringVar(&watchExtensions, "watch.extensions", ".go,.css,.js", "comma separated extensions to watch for recompile")
}

type Watcher interface {
	Watch(reload []chan bool)
}

type watcher struct{}

func (w *watcher) Watch(reloadCh []chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading paths: %v\n", err)
		return
	}

	d := newDebounce()
	defer d.timer.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if !slices.Contains(strings.Split(watchExtensions, ","), filepath.Ext(event.Name)) {
				continue
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) ||
				event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {

				d.Trigger(reloadCh)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

func newDebounce() *debounce {
	delay := 100 * time.Millisecond

	return &debounce{
		timer: time.NewTimer(delay),
		delay: delay,
	}
}

type debounce struct {
	timer     *time.Timer
	delay     time.Duration
	lastEvent time.Time
}

func (d *debounce) Trigger(reloadCh []chan bool) {
	now := time.Now()
	if now.Sub(d.lastEvent) > d.delay {
		for _, ch := range reloadCh {
			select {
			case ch <- true:
			default:
			}
		}

		d.lastEvent = now
		return
	}

	if !d.timer.Stop() {
		<-d.timer.C
	}

	d.timer.Reset(d.delay)
}
