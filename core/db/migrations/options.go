package migrations

import "fmt"

// migrationOptions are the options for the migration
type Option func() error

// UseMigrationFolder sets the folder for migrations
func UseMigrationFolder(f string) Option {
	return func() error {
		folder = f

		return nil
	}
}

// Applies migration options like setting the folder
// for migrations.
func Apply(options ...Option) error {
	for _, option := range options {
		if err := option(); err != nil {
			return fmt.Errorf("error applying migration option: %w", err)
		}
	}

	return nil
}
