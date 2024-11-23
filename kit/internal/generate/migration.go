package generate

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

var (
	// migrationTemplate is the template for generating migrations
	//go:embed migration.sql.tmpl
	migrationsTemplate string
)

// Migration generator function
func Migration(migrationsFolder, name string) error {
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf(
		"%s_%s.sql",
		timestamp,
		name,
	)

	t, err := template.New("migration").Parse(migrationsTemplate)
	if err != nil {
		return fmt.Errorf("error parsing migrations template: %w", err)
	}

	// Destination file name
	fileName = filepath.Join(migrationsFolder, fileName)
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
