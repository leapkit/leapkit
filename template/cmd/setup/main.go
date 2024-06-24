package main

import (
	"fmt"

	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/template/internal"
	"github.com/leapkit/leapkit/template/internal/migrations"
	"github.com/paganotoni/tailo"

	// Load environment variables
	_ "github.com/leapkit/leapkit/core/tools/envload"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Setup tailo to compile tailwind css.
	err := tailo.Setup()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("✅ Tailwind CSS setup successfully")
	err = db.Create(internal.DatabaseURL)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("✅ Database created successfully")
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
}
