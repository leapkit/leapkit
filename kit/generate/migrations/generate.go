package migrations

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

var (
	// migrationsFolder is the base folder for migrations
	folder = filepath.Join(
		"internal", "app", "database", "migrations",
	)

	// migrationTemplate is the template for generating migrations
	templ = `-- {{.Timestamp}} - {{.Name }} migration`
)

func Folder() string {
	return folder
}

// Returns the template for migrations.
func Template() string {
	return templ
}

func generate(m instance, opts ...Option) error {
	// Applying options before generating the migration
	Apply(opts...)

	t, err := template.New("migration").Parse(Template())
	if err != nil {
		return fmt.Errorf("error parsing migrations template: %w", err)
	}

	// Destination file name
	name := filepath.Join(Folder(), m.Filename())
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("error creating migration file: %w", err)
	}

	err = t.ExecuteTemplate(f, "migration", m)
	if err != nil {
		return fmt.Errorf("error executing migrations template: %w", err)
	}

	fmt.Printf("✅ Migration file `%v` generated\n", name)
	return nil
}
