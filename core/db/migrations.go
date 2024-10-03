package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// RunMigrationsDir receives a folder and a database URL
// to apply the migrations to the database.
func RunMigrationsDir(dir string, conn *sql.DB) error {
	migrator := NewMigrator(conn)
	err := migrator.Setup()
	if err != nil {
		return fmt.Errorf("error setting up migrations: %w", err)
	}

	exp := regexp.MustCompile("(\\d{14})_(.*).sql")
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking migrations directory: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		matches := exp.FindStringSubmatch(path)
		if len(matches) != 3 {
			return nil
		}

		name := matches[2]
		timestamp := matches[1]
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error opening migration file: %w", err)
		}

		err = migrator.Run(timestamp, name, string(content))
		if err != nil {
			return fmt.Errorf("error running migration %s: %w", path, err)
		}

		return nil
	})
}

// RunMigrations by checking in the migrations database
// table, each of the adapters take care of this.
func RunMigrations(fs embed.FS, conn *sql.DB) error {
	dir, err := fs.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}

	migrator := NewMigrator(conn)
	err = migrator.Setup()
	if err != nil {
		return fmt.Errorf("error setting up migrations: %w", err)
	}

	exp := regexp.MustCompile("(\\d{14})_(.*).sql")
	for _, v := range dir {
		if v.IsDir() {
			continue
		}

		matches := exp.FindStringSubmatch(v.Name())
		if len(matches) != 3 {
			continue
		}

		timestamp := matches[1]
		name := matches[2]
		content, err := fs.ReadFile(v.Name())
		if err != nil {
			return fmt.Errorf("error opening migration file: %w", err)
		}

		err = migrator.Run(timestamp, name, string(content))
		if err != nil {
			return fmt.Errorf("error running migration: %w", err)
		}
	}

	return nil
}
