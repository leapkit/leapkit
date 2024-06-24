package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/leapkit/leapkit/kit/internal/assets"
	"github.com/leapkit/leapkit/kit/internal/rebuilder"
	"github.com/paganotoni/tailo"
)

var (
	assetsInput     string
	assetsOutput    string
	watchExtensions string

	tailoInput  string
	tailoOutput string
	tailoConfig string
)

func init() {
	flag.StringVar(&assetsInput, "assets.input", "internal/assets", "the input folder for the assets")
	flag.StringVar(&assetsOutput, "assets.output", "public", "the output folder for the assets")
	flag.StringVar(&tailoInput, "tailo.input", "internal/assets/application.css", "")
	flag.StringVar(&tailoOutput, "tailo.output", "public/application.css", "")
	flag.StringVar(&tailoConfig, "tailo.config", "tailwind.config.js", "")

	flag.StringVar(&watchExtensions, "watch.extensions", ".go,.css,.js", "comma separated extensions to watch for recompile")
}

func serve(_ []string) error {
	err := rebuilder.Start(
		"cmd/app/main.go",

		// Run the tailo watcher so when changes are made to
		// the html code it rebuilds css.
		rebuilder.WithRunner(tailo.WatcherFn(
			tailo.UseInputPath(tailoInput),
			tailo.UseOutputPath(tailoOutput),
			tailo.UseConfigPath(tailoConfig),
		)),

		// Run the assets watcher.
		rebuilder.WithRunner(assets.Watch(assetsInput, assetsOutput)),
		rebuilder.WatchExtension(strings.Split(watchExtensions, ",")...),
	)

	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}

	return err
}
