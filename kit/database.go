package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leapkit/core/db"

	// Loading .env file
	_ "github.com/leapkit/core/tools/envload"

	// Postgres driver
	_ "github.com/lib/pq"

	// Sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// database provides operations to manage the database
// during development. It can create, drop and run migrations.
func database(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: database <command>")

		return nil
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		fmt.Println(" DATABASE_URL is not set")

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
			return err
		}

		err = db.RunMigrationsDir(filepath.Join("internal", "migrations"), conn)
		if err != nil {
			return err
		}

		fmt.Println("✅ Migrations ran successfully")
	case "create":
		err := db.Create(url)
		if err != nil {
			return err
		}

		fmt.Println("✅ Database created successfully")

	case "drop":
		err := db.Drop(url)
		if err != nil {
			return err
		}

		fmt.Println("✅ Database dropped successfully")
	default:
		fmt.Println("command not found")

		return nil
	}

	return nil
}
