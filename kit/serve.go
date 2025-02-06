package main

import (
	"fmt"
	"strings"

	"github.com/leapkit/leapkit/kit/internal/rebuilder"
	flag "github.com/spf13/pflag"
)

var watchExtensions string

func init() {
	flag.StringVar(&watchExtensions, "watch.extensions", ".go,.css,.js", "comma separated extensions to watch for recompile")
}

func serve(_ []string) error {
	exts := rebuilder.WatchExtension(strings.Split(watchExtensions, ",")...)
	err := rebuilder.Start("cmd/app/main.go", exts)
	if err != nil {
		fmt.Println("[error] starting the server:", err)
	}

	return err
}
