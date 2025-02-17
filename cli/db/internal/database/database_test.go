package database_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/leapkit/leapkit/cli/db/internal/database"
)

func TestMigrate(t *testing.T) {
	testCases := []struct {
		name            string
		url             string
		output          string
		migrationFolder bool
	}{
		{
			name:            "No migrations folder",
			url:             "file::memory:?cache=shared",
			output:          "[error] error running migrations: error walking migrations directory: lstat internal/migrations: no such file or directory\n",
			migrationFolder: false,
		},

		{
			name:            "Run Successfully",
			url:             "file::memory:?cache=shared",
			output:          "✅ Migrations ran successfully\n",
			migrationFolder: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.RemoveAll("internal")
			// Set the DATABASE_URL to the temporary file
			os.Setenv("DATABASE_URL", tc.url)
			if tc.migrationFolder {
				// Create temporary migrations folder in internal/migrations
				err := os.MkdirAll("internal/migrations", 0755)
				if err != nil {
					t.Fatalf("error creating migrations folder: %v", err)
				}

				defer os.RemoveAll("internal")
			}

			stdout := os.Stdout
			stderr := os.Stderr

			f, _ := os.Create("output")
			defer os.Remove("output")

			os.Stdout = f
			os.Stderr = f

			os.Args = []string{"db", "migrate"}
			// main.go call
			err := database.Exec()
			if err != nil {
				fmt.Printf("[error] %v\n", err)
			}

			os.Stdout = stdout
			os.Stderr = stderr

			f.Close()
			f, _ = os.Open("output")

			out, _ := io.ReadAll(f)
			if !strings.Contains(string(out), tc.output) {
				t.Fatalf("unexpected output: %v", string(out))
			}
		})
	}

}

func TestCreate(t *testing.T) {
	t.Run("Create SQLite", func(t *testing.T) {
		// Set the DATABASE_URL to the temporary file
		os.Setenv("DATABASE_URL", "test.db")
		defer os.Remove("test.db")

		stdout := os.Stdout
		stderr := os.Stderr

		f, _ := os.Create("output")
		defer os.Remove("output")

		os.Stdout = f
		os.Stderr = f

		os.Args = []string{"db", "create"}
		// main.go call
		err := database.Exec()
		if err != nil {
			fmt.Printf("[error] %v\n", err)
		}

		os.Stdout = stdout
		os.Stderr = stderr

		f.Close()
		f, _ = os.Open("output")

		out, _ := io.ReadAll(f)
		if string(out) != "✅ Database created successfully\n" {
			t.Fatalf("unexpected output: %v", string(out))
		}

		// Check if the file exists
		_, err = os.Stat("test.db")
		if err != nil || os.IsNotExist(err) {
			t.Fatalf("file does not exist")
		}
	})
}

func TestDrop(t *testing.T) {
	t.Run("Drop SQLite", func(t *testing.T) {
		// Set the DATABASE_URL to the temporary file
		os.Setenv("DATABASE_URL", "test.db")
		defer os.Remove("test.db")

		stdout := os.Stdout
		stderr := os.Stderr

		f, _ := os.Create("output")
		defer os.Remove("output")

		os.Stdout = f
		os.Stderr = f

		os.Args = []string{"db", "drop"}
		// main.go call
		err := database.Exec()
		if err != nil {
			fmt.Printf("[error] %v\n", err)
		}

		os.Stdout = stdout
		os.Stderr = stderr

		f.Close()
		f, _ = os.Open("output")

		out, _ := io.ReadAll(f)
		if string(out) != "✅ Database dropped successfully\n" {
			t.Fatalf("unexpected output: %v", string(out))
		}

		// Check if the file exists
		_, err = os.Stat("test.db")
		if err == nil || !os.IsNotExist(err) {
			t.Fatalf("file exists")
		}

	})

}

func TestReset(t *testing.T) {

	testCases := []struct {
		name            string
		url             string
		output          string
		migrationFolder bool
	}{

		{
			name:            "No migrations folder",
			url:             "test.db",
			output:          "[error] error running migrations: error walking migrations directory: lstat internal/migrations: no such file or directory\n",
			migrationFolder: false,
		},

		{
			name:            "Run Successfully",
			url:             "test.db",
			output:          "✅ Database reset successfully\n",
			migrationFolder: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the DATABASE_URL to the temporary file
			os.Setenv("DATABASE_URL", tc.url)
			defer os.Remove(tc.url)

			if tc.migrationFolder {
				// Create temporary migrations folder in internal/migrations
				err := os.MkdirAll("internal/migrations", 0755)
				if err != nil {
					t.Fatalf("error creating migrations folder: %v", err)
				}

				defer os.RemoveAll("internal")
			}

			stdout := os.Stdout
			stderr := os.Stderr

			f, _ := os.Create("output")
			defer os.Remove("output")

			os.Stdout = f
			os.Stderr = f

			os.Args = []string{"db", "reset"}
			// main.go call
			err := database.Exec()
			if err != nil {
				fmt.Printf("[error] %v\n", err)
			}

			os.Stdout = stdout
			os.Stderr = stderr

			f.Close()
			f, _ = os.Open("output")

			out, _ := io.ReadAll(f)
			if string(out) != tc.output {
				t.Fatalf("unexpected output: %v", string(out))
			}
		})
	}

}
