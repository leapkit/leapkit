package main_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestMigrate(t *testing.T) {

	testCases := []struct {
		name            string
		url             string
		output          string
		migrationFolder bool
	}{
		{
			name:            "No DATABASE_URL",
			url:             "",
			output:          "[error] DATABASE_URL is not set\n",
			migrationFolder: true,
		},

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
			// Set the DATABASE_URL to the temporary file
			os.Setenv("DATABASE_URL", tc.url)
			if tc.migrationFolder {
				// Create temporary migrations folder in internal/migrations
				err := os.MkdirAll("internal/migrations", 0755)
				if err != nil {
					t.Fatalf("error creating migrations folder: %v", err)
				}

				defer os.RemoveAll("internal/migrations")
			}

			cmd := exec.Command("go", "run", ".", "db", "migrate")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running command: %v", err)
			}

			if string(out) != tc.output {
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
		cmd := exec.Command("go", "run", ".", "db", "create")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error running command: %v", err)
		}

		if string(out) != "✅ Database created successfully\n" {
			t.Fatalf("unexpected output: %v", string(out))
		}

		// Check if the file exists
		_, err = os.Stat("test.db")
		if err != nil {
			t.Fatalf("error checking file")
		}

		if os.IsNotExist(err) {
			t.Fatalf("file does not exist")
		}
	})

}

func TestDrop(t *testing.T) {
	t.Run("Drop SQLite", func(t *testing.T) {
		// Set the DATABASE_URL to the temporary file
		os.Setenv("DATABASE_URL", "test.db")
		defer os.Remove("test.db")
		cmd := exec.Command("go", "run", ".", "db", "drop")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("error running command: %v", err)
		}

		if string(out) != "✅ Database dropped successfully\n" {
			t.Fatalf("unexpected output: %v", string(out))
		}

		// Check if the file exists
		_, err = os.Stat("test.db")
		if err == nil {
			t.Fatalf("file exists")
		}

		if !os.IsNotExist(err) {
			t.Fatalf("error checking file:")
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
			name:            "No DATABASE_URL",
			url:             "",
			output:          "[error] DATABASE_URL is not set\n",
			migrationFolder: true,
		},

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

				defer os.RemoveAll("internal/migrations")
			}

			cmd := exec.Command("go", "run", ".", "db", "reset")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running command: %v", err)
			}

			if string(out) != tc.output {
				t.Fatalf("unexpected output: %v", string(out))
			}
		})
	}

}
