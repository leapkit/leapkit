package main

import (
	"cmp"
	"fmt"
	"os"
	"os/exec"

	"github.com/leapkit/leapkit/template/internal"
	"github.com/leapkit/leapkit/template/internal/migrations"
	"go.leapkit.dev/core/db"

	// Load environment variables
	_ "go.leapkit.dev/core/tools/envload"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Running the tailo setup command
	cmd := exec.Command("go", "run", "github.com/paganotoni/tailo/cmd/build@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("✅ Tailwind CSS setup successfully")
	err = db.Create(cmp.Or(os.Getenv("DATABASE_URL"), "database.db?_timeout=5000&_sync=1"))
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
