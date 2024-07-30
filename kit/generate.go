package main

import (
	"fmt"

	"github.com/leapkit/leapkit/kit/generate"
)

func generateWith(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: generate <generator_name>")
		fmt.Println("Available commands:")
		fmt.Println("  - migration [name]")
		fmt.Println("  - action [action|folder/action]")
		fmt.Println("")
		return nil
	}

	switch args[1] {
	case "migration":
		if len(args) < 3 {
			fmt.Println("Usage: generate migration <name>")

			return nil
		}

		err := generate.Migration(args[2])
		if err != nil {
			return err
		}
	case "action":
		usage := func() error {
			fmt.Println("Usage: generate action [action|folder/action]")
			return nil
		}
		if len(args) < 3 {
			return usage()
		}

		if args[2] == "" {
			return usage()
		}

		err := generate.Action(args[2])
		if err != nil {
			return err
		}
	default:
		fmt.Println("Usage: generate [generator]")
		return nil
	}

	return nil
}
