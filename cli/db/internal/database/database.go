package database

import (
	"cmp"
	"fmt"
	"os"

	"github.com/leapkit/leapkit/core/db"
	flag "github.com/spf13/pflag"

	// Loading .env file
	_ "github.com/leapkit/leapkit/core/tools/envload"

	// Postgres driver
	_ "github.com/lib/pq"

	// Sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Exec provides operations to manage the database
// during development. It can create, drop and run migrations.
func Exec() error {
	flag.Parse()

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: database <command>")

		return nil
	}

	url := cmp.Or(os.Getenv("DATABASE_URL"), "database.db?_timeout=5000&_sync=1")

	switch args[1] {
	case "migrate":
		err := runMigrations(url)
		if err != nil {
			fmt.Printf("[error] %v\n", err)
		}

		fmt.Println("✅ Migrations ran successfully")
	case "create":
		err := db.Create(url)
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}

		fmt.Println("✅ Database created successfully")

	case "drop":
		err := db.Drop(url)
		if err != nil {
			return fmt.Errorf("error dropping database: %w", err)
		}

		fmt.Println("✅ Database dropped successfully")

	case "reset":
		err := db.Drop(url)
		if err != nil {
			return fmt.Errorf("error dropping database: %w", err)
		}

		err = db.Create(url)
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}

		if err := runMigrations(url); err != nil {
			return err
		}

		fmt.Println("✅ Database reset successfully")

	case "generate_migration":
		if len(args) < 3 {
			fmt.Println("Usage: database generate_migration <migration_name>")
			return nil
		}

		err := newMigration(args[2])
		if err != nil {
			return err
		}
	default:
		fmt.Println("command not found")

		return nil
	}

	return nil
}
