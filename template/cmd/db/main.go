package main

import (
	"fmt"
	"os"

	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/template/internal"
	"github.com/leapkit/leapkit/template/internal/migrations"

	// Load environment variables
	_ "github.com/leapkit/leapkit/core/tools/envload"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tools database <command>")
		fmt.Println("Available commands:")
		fmt.Println(" - migrate")
		fmt.Println(" - create")
		fmt.Println(" - drop")

		return
	}

	switch os.Args[1] {
	case "migrate":
		conn, err := internal.DB()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = db.RunMigrations(migrations.All, conn)
		if err != nil {
			fmt.Println(err)

			return
		}

		fmt.Println("✅ Migrations ran successfully")
	case "create":
		err := db.Create(internal.DatabaseURL)
		if err != nil {
			fmt.Println(err)

			return
		}

		fmt.Println("✅ Database created successfully")

	case "drop":
		err := db.Drop(internal.DatabaseURL)
		if err != nil {
			fmt.Println(err)

			return
		}

		fmt.Println("✅ Database dropped successfully")

	default:
		fmt.Println("command not found")

		return
	}
}
