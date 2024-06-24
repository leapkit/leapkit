package migrations

// Migrator that can take care of migration operations
// depending on the database.
type Migrator interface {
	// Setup the database, p.e. creating the migrations table.
	Setup() error
	// Run specific migration timestamp and SQL.
	Run(timestamp, sql string) error
}
