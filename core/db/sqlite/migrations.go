package sqlite

import "fmt"

// Setup the sqlite database to be ready to have the migrations inside.
func (a *adapter) Setup() error {
	_, err := a.conn.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (timestamp TEXT);")
	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	return nil
}

// Run a particular database migration and inserting its timestamp
// on the migrations table.
func (a *adapter) Run(timestamp, sql string) error {
	var exists bool
	row := a.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE timestamp = $1)", timestamp)
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("error running migration: %w", err)
	}

	if !exists {
		_, err = a.conn.Exec(sql)
		if err != nil {
			return fmt.Errorf("error running migration: %w", err)
		}

		_, err = a.conn.Exec("INSERT INTO schema_migrations (timestamp) VALUES ($1);", timestamp)
		if err != nil {
			return fmt.Errorf("error running migration: %w", err)
		}

		fmt.Printf("✅ Migration %v applied\n", timestamp)
	}

	return nil
}
