package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leapkit/leapkit/core/db"
	flag "github.com/spf13/pflag"

	// Loading .env file
	_ "github.com/leapkit/leapkit/core/tools/envload"

	// Postgres driver
	_ "github.com/lib/pq"

	// Sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	// migrationsFolder is the folder where the migrations are stored
	migrationsFolder string
)

func init() {
	flag.StringVar(&migrationsFolder, "migrations.folder", filepath.Join("internal", "migrations"), "the folder where the migrations are stored")
}

// database provides operations to manage the database
// during development. It can create, drop and run migrations.
func database(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: database <command>")

		return nil
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		fmt.Println("[error] DATABASE_URL is not set")

		return nil
	}

	switch args[1] {
	case "migrate":
		driver := "sqlite3"
		if strings.HasPrefix(url, "postgres") {
			driver = "postgres"
		}

		conn, err := sql.Open(driver, url)
		if err != nil {
			return fmt.Errorf("error opening connection: %w", err)
		}

		err = db.RunMigrationsDir(migrationsFolder, conn)
		if err != nil {
			return fmt.Errorf("error running migrations: %w", err)
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

		driver := "sqlite3"
		if strings.HasPrefix(url, "postgres") {
			driver = "postgres"
		}

		conn, err := sql.Open(driver, url)
		if err != nil {
			return fmt.Errorf("error opening connection: %w", err)
		}

		err = db.RunMigrationsDir(migrationsFolder, conn)
		if err != nil {
			return fmt.Errorf("error running migrations: %w", err)
		}

		fmt.Println("✅ Database reset successfully")
	default:
		fmt.Println("command not found")

		return nil
	}

	return nil
}
