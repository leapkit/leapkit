package migrations

import "path/filepath"

var (
	// migrationsFolder is the base folder for migrations
	folder = filepath.Join(
		"internal", "app", "database", "migrations",
	)

	// migrationTemplate is the template for generating migrations
	template = `-- {{.Timestamp}} - {{.Name }} migration`
)

// Returns the folder for migrations
func Folder() string {
	return folder
}

// Returns the template for migrations.
func Template() string {
	return template
}
