package main

import (
	"fmt"

	"github.com/leapkit/leapkit/core/tools/rebuilder"
	"github.com/leapkit/leapkit/template/internal"
	"github.com/paganotoni/tailo"

	// Load environment variables
	_ "github.com/leapkit/leapkit/core/tools/envload"
)

func main() {
	err := rebuilder.Start(
		"cmd/app/main.go",

		// Run the tailo watcher so when changes are made to
		// the html code it rebuilds css.
		rebuilder.WithRunner(tailo.WatcherFn(
			tailo.UseInputPath("internal/assets/application.css"),
			tailo.UseOutputPath("public/application.css"),
			tailo.UseConfigPath("tailwind.config.js"),
		)),

		// Run the assets watcher.
		rebuilder.WithRunner(internal.Assets.Watch),
		rebuilder.WatchExtension(".go", ".css", ".js"),
	)

	if err != nil {
		fmt.Println("error starting the dev", err)
	}
}
