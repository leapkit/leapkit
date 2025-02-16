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
		fmt.Println("  - version [command]")
		fmt.Println("  - update")
		fmt.Println("")

		fmt.Println("Available flags:")
		flag.PrintDefaults()

		return
	}

	var err error
	switch os.Args[1] {
	case "new":
		newmodule(os.Args[1:])
	case "update":
		err = update()
	case "version", "v":
		version(os.Args[1:])
	default:
		fmt.Printf("Unknown command `%v`.\n\n", os.Args[1])
		return
	}

	if err != nil {
		fmt.Printf("[error] %v\n", err)
	}
}
