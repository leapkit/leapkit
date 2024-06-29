package main

import (
	"fmt"
	"runtime/debug"
)

func version(_ []string) {
	version := "(main)"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}

	fmt.Printf("Kit version: %v\n", version)
}
