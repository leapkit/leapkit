package database

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/leapkit/leapkit/core/db"
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

func runMigrations(url string) error {
	driver := "sqlite3"
	if strings.HasPrefix(url, "postgres") {
		driver = "postgres"
	}

	conn, err := sql.Open(driver, url)
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	err = db.RunMigrationsDir(migrationFolder, conn)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}
