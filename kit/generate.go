package main

import (
	"fmt"
	"path/filepath"

	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/core/db/migrations"
)

func generate(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: generate <generator_name>")

		return nil
	}

	if args[1] != "migration" {
		fmt.Println("Usage: generate [generator]")

		return nil
	}

	if len(args) < 3 {
		fmt.Println("Usage: generate migration <name>")

		return nil
	}

	err := db.GenerateMigration(
		args[2], // name of the migration

		// This is the path to the migrations folder
		migrations.UseMigrationFolder(filepath.Join("internal", "migrations")),
	)

	if err != nil {
		return err
	}

	return nil
}
