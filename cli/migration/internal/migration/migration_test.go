package migration_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/leapkit/leapkit/cli/migration/internal/migration"
	"github.com/spf13/pflag"
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
		os.Args = []string{"migration", "new", "create_users_table"}
		err = migration.Migration()
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

		customMigrationFolder := "internal/database/migrations"

		flagValue := stringValue(customMigrationFolder)
		pflag.CommandLine.AddFlag(&pflag.Flag{
			Name:      "migration.folder",
			Shorthand: "f",
			Value:     &flagValue,
			DefValue:  customMigrationFolder,
			Usage:     "test",
		})

		os.Args = []string{"migration", "new", "create_users_table", fmt.Sprintf("--migration.folder=%s", customMigrationFolder)}

		// Create a new migration
		err = migration.Migration()
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

	t.Run("correct incomplete command", func(t *testing.T) {
		current := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		defer func() {
			os.Stdout = current
		}()

		os.Args = []string{"migration", "foo"}
		err := migration.Migration()
		if err != nil {
			t.Fatalf("error creating migration: %v", err)
		}

		w.Close()
		out, _ := io.ReadAll(r)
		fmt.Println(string(out))

		if !bytes.Contains(out, []byte("Usage: migration <command>")) {
			t.Errorf("Expected 'Usage: migration <command>', got: %v", string(out))
		}
	})

	t.Run("correct invalida command", func(t *testing.T) {
		current := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		defer func() {
			os.Stdout = current
		}()

		os.Args = []string{"migration", "invalid", "value"}
		err := migration.Migration()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		w.Close()
		out, _ := io.ReadAll(r)
		fmt.Println(string(out))

		if !bytes.Contains(out, []byte("command not found")) {
			t.Errorf("Expected 'command not found', got: %v", string(out))
		}
	})
}

type stringValue string

func (s *stringValue) String() string { return string(*s) }
func (s *stringValue) Type() string   { return "string" }
func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
