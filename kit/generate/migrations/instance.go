package migrations

import (
	"fmt"
	"path/filepath"
	"time"
)

// instance of a migration file, it has a name and a
// timestamp.
type instance struct {
	Name      string
	Timestamp string
}

// Filename of the migration.
func (m instance) Filename() string {
	return fmt.Sprintf("%s_%s.sql", m.Timestamp, m.Name)
}

// Generates a new migration with the passed name and
// applies current timestamp to it.
func New(name string) instance {
	return instance{
		Name:      name,
		Timestamp: time.Now().Format("20060102150405"),
	}
}

// GenerateMigration in the migrations folder using the migrations template
func (m instance) Generate() error {
	return generate(
		m, // instance of the migration
		// This is the path to the migrations folder
		UseFolder(filepath.Join("internal", "migrations")),
	)
}
