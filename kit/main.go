package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func main() {
	// Parse flags
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage: kit <command>")
		fmt.Println("Available commands:")
		fmt.Println("  - new [name]")
		fmt.Println("  - database [command]")
		fmt.Println("  - generate [generator]")
		fmt.Println("  - serve [command]")
		fmt.Println("  - version [command]")
		fmt.Println("")

		fmt.Println("Available flags:")
		flag.PrintDefaults()

		return
	}

	switch os.Args[1] {
	case "new":
		newmodule(os.Args[1:])
	case "serve", "s", "dev":
		serve(os.Args[1:])
	case "database", "db":
		database(os.Args[1:])
	case "generate", "gen", "g":
		generate(os.Args[1:])
	case "version", "v":
		version(os.Args[1:])
	default:
		fmt.Printf("Unknown command `%v`.\n\n", os.Args[1])
	}
}
