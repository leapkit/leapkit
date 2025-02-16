package migration

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
)

var (
	// migrationTemplate is the template for generating migrations
	//go:embed migration.sql.tmpl
	migrationTemplate string

	// migrationFolder is the folder where the migrations are stored
	migrationFolder string
)

func init() {
	flag.StringVar(&migrationFolder, "migration.folder", filepath.Join("internal", "migrations"), "the folder where the migrations are stored")
}

func Migration() error {
	pflag.Parse()

	if f := pflag.Lookup("migration.folder"); f != nil {
		migrationFolder = f.Value.String()
	}

	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: migration <command>")

		return nil
	}

	switch args[1] {
	case "new":
		err := newMigration(args[2])
		if err != nil {
			return fmt.Errorf("error creating migration: %w", err)
		}
	default:
		fmt.Println("command not found")
	}

	return nil
}

// newMigration generator function
func newMigration(name string) error {
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf(
		"%s_%s.sql",
		timestamp,
		name,
	)

	t, err := template.New("migration").Parse(migrationTemplate)
	if err != nil {
		return fmt.Errorf("error parsing migrations template: %w", err)
	}

	// Destination file name
	fileName = filepath.Join(migrationFolder, fileName)
	err = os.MkdirAll(filepath.Dir(fileName), 0700)
	if err != nil {
		return fmt.Errorf("error creating migrations folder: %w", err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating migration file: %w", err)
	}

	err = t.ExecuteTemplate(f, "migration", map[string]string{
		"Name":      name,
		"Timestamp": timestamp,
	})

	if err != nil {
		return fmt.Errorf("error executing migrations template: %w", err)
	}

	fmt.Printf("✅ Migration file `%v` generated\n", name)
	return nil
}
