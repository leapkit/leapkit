package main

import (
	"fmt"

	"github.com/leapkit/leapkit/kit/internal/assets"
	"github.com/leapkit/leapkit/kit/internal/rebuilder"
	"github.com/paganotoni/tailo"
)

func serve(_ []string) error {
	// TODO: need to be flags
	var (
		inputFolder  = "internal/assets"
		outputFolder = "public"
		extensions   = []string{".go", ".css", ".js"}

		//Tailo
		tailwindInput  = "internal/assets/application.css"
		tailwindOutput = "public/application.css"
		tailwindConfig = "tailwind.config.js"
	)

	err := rebuilder.Start(
		"cmd/app/main.go",

		// Run the tailo watcher so when changes are made to
		// the html code it rebuilds css.
		rebuilder.WithRunner(tailo.WatcherFn(
			tailo.UseInputPath(tailwindInput),
			tailo.UseOutputPath(tailwindOutput),
			tailo.UseConfigPath(tailwindConfig),
		)),

		// Run the assets watcher.
		rebuilder.WithRunner(assets.Watch(inputFolder, outputFolder)),
		rebuilder.WatchExtension(extensions...),
	)

	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}

	return err
}
