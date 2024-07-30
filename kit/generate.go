package main

import (
	"fmt"
	"strings"

	gen "github.com/leapkit/leapkit/kit/generate"
)

func generate(args []string) error {
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

		err := gen.New(gen.Params{
			Kind: "migration",
			Name: args[2],
		}).Generate()

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

		path := strings.Split(args[2], "/")
		err := gen.New(gen.Params{
			Kind: "action",
			Path: path,
		}).Generate()

		if err != nil {
			return err
		}
	default:
		fmt.Println("Usage: generate [generator]")
		return nil
	}

	return nil
}
