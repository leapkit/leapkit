package generate_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/leapkit/leapkit/kit/internal/generate"
)

func TestGenerateMigration(t *testing.T) {
	bd, _ := os.Getwd()
	defer os.Chdir(bd)

	t.Run("correct generate migration", func(t *testing.T) {
		wd := t.TempDir()
		err := os.Chdir(wd)
		if err != nil {
			t.Fatalf("error changing directory: %v", err)
		}

		migrationFolder := "internal/migrations"

		// Create a new migration
		err = generate.Migration(migrationFolder, "create_users_table")
		if err != nil {
			t.Fatalf("error creating migration: %v", err)
		}

		var migrationPath string
		// Check if the migration file was created
		filepath.Walk(filepath.Join(wd, migrationFolder), func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".sql" {
				migrationPath = path
				return nil
			}

			return nil
		})

		if migrationPath == "" {
			t.Fatal("migration file not created")
		}

		migrationPath, err = filepath.Rel(wd, migrationPath)
		if err != nil {
			t.Fatalf("error getting relative path: %v", err)
		}

		// Check if the migration file is not empty
		file, err := os.Open(migrationPath)
		if err != nil {
			t.Fatalf("error opening migration file: %v", err)
		}

		defer file.Close()

		// read the file content
		bc, err := os.ReadFile(migrationPath)
		if err != nil {
			t.Fatalf("error reading migration file: %v", err)
		}

		if bytes.Contains(bc, []byte(migrationPath)) {
			t.Fatalf("migration should not contain the full path")
		}
	})

	t.Run("correct generate migration in custom folder", func(t *testing.T) {
		wd := t.TempDir()
		err := os.Chdir(wd)
		if err != nil {
			t.Fatalf("error changing directory: %v", err)
		}

		customMigrationFolder := "internal/migrations"
		// Create a new migration
		err = generate.Migration(customMigrationFolder, "create_users_table")
		if err != nil {
			t.Fatalf("error creating migration: %v", err)
		}

		var migrationPath string
		// Check if the migration file was created
		filepath.Walk(filepath.Join(wd, customMigrationFolder), func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".sql" {
				migrationPath = path
				return nil
			}

			return nil
		})

		if migrationPath == "" {
			t.Fatal("migration file not created")
		}

		migrationPath, err = filepath.Rel(wd, migrationPath)
		if err != nil {
			t.Fatalf("error getting relative path: %v", err)
		}

		// Check if the migration file is not empty
		file, err := os.Open(migrationPath)
		if err != nil {
			t.Fatalf("error opening migration file: %v", err)
		}

		defer file.Close()

		// read the file content
		bc, err := os.ReadFile(migrationPath)
		if err != nil {
			t.Fatalf("error reading migration file: %v", err)
		}

		if bytes.Contains(bc, []byte(migrationPath)) {
			t.Fatalf("migration should not contain the full path")
		}
	})
}
