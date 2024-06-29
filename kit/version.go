package main

import (
	"fmt"
	"runtime/debug"
)

// version prints the version of the kit by taking it
// from the build info or setting it to "(main)" if
// not available.
func version(_ []string) {
	version := "(main)"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}

	fmt.Printf("Kit version: %v\n", version)
}
