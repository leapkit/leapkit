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
func (a *adapter) Run(timestamp, name, sql string) error {
	var exists bool
	row := a.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE timestamp = $1)", timestamp)
	err := row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("error running migration: %w", err)
	}

	tx, err := a.conn.Begin()
	if err != nil {
		return fmt.Errorf("error running migration: %w", err)
	}

	defer func() {
		// If there is an error, rollback the transaction.
		if err != nil {
			fmt.Printf("🚨 Rolling back migration: %s\n", err.Error())
			tx.Rollback()
		}
	}()

	if exists {
		return nil
	}

	_, err = tx.Exec(sql)
	if err != nil {
		err = fmt.Errorf("error running migration: %w", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO schema_migrations (timestamp) VALUES ($1);", timestamp)
	if err != nil {
		err = fmt.Errorf("error running migration: %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error running migration: %w", err)
	}

	fmt.Printf("✅ Migration %v (%v) applied.\n", name, timestamp)

	return nil
}
