package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

func main() {
	// Parse flags
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage: kit <command>")
		fmt.Println("Available commands:")
		fmt.Println(" - database [command]")
		fmt.Println(" - generate [generator]")

		return
	}

	switch os.Args[1] {
	case "serve", "s", "dev":
		serve(os.Args[1:])
	case "database", "db":
		database(os.Args[1:])
	case "generate", "gen", "g":
		generate(os.Args[1:])
	case "version", "v":
		version := "(main)"
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
		}

		fmt.Printf("Kit version: %v\n", version)
	default:
		fmt.Printf("Unknown command `%v`.\n\n", os.Args[1])
	}
}
