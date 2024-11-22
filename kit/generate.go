package main

import (
	"fmt"

	"github.com/leapkit/leapkit/kit/internal/generate"
)

func generateWith(args []string) error {
	mainUsage := func() {
		fmt.Println("Usage: generate <generator_name>")
		fmt.Println("Available commands:")
		fmt.Println("  - migration [name]")
		fmt.Println("  - action [action|folder/action]")
		fmt.Println("  - handler [name|folder/name]")
		fmt.Println("")
	}

	if len(args) < 2 {
		mainUsage()
		return nil
	}

	switch args[1] {
	case "migration":
		if len(args) < 3 {
			fmt.Println("Usage: generate migration <name>")

			return nil
		}

		err := generate.Migration(migrationsFolder, args[2])
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

	case "handler":
		usage := func() error {
			fmt.Println("Usage: generate handler [name|folder/name]")
			return nil
		}
		if len(args) < 3 {
			return usage()
		}

		if args[2] == "" {
			return usage()
		}

		err := generate.Handler(args[2])
		if err != nil {
			return err
		}

	default:
		mainUsage()
		return nil
	}

	return nil
}
