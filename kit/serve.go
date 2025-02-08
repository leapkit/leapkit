package main

import (
	"fmt"
	"strings"

	"github.com/leapkit/leapkit/kit/internal/rebuilder"
	flag "github.com/spf13/pflag"
)

var watchExtensions string

func init() {
	flag.StringVar(&watchExtensions, "watch.extensions", ".go", "comma separated extensions to watch for recompile")
}

func serve(_ []string) error {
	err := rebuilder.Start(
		"cmd/app/main.go",

		// expecifying the extensions to watch for
		rebuilder.WatchExtension(
			strings.Split(watchExtensions, ",")...,
		),
	)
	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}

	return err
}
