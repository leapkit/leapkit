package main_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGenerateMigration(t *testing.T) {
	testCases := []struct {
		name            string
		output          string
		migrationFolder bool
	}{
		{
			name:            "No migrations folder",
			output:          "[error] error creating migration file:",
			migrationFolder: false,
		},
		{
			name:            "Migration folder",
			output:          "✅ Migration file `internal/migrations/",
			migrationFolder: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.migrationFolder {
				// Create temporary migrations folder in internal/migrations
				err := os.MkdirAll("internal/migrations", 0755)
				if err != nil {
					t.Fatalf("error creating migrations folder: %v", err)
				}

				defer os.RemoveAll("internal/migrations")
			}

			cmd := exec.Command("go", "run", ".", "gen", "migration", "create_users_table")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running command: %v", err)
			}

			expected := strings.Replace(tc.output, " ", "", -1)
			ouput := strings.Replace(string(out), " ", "", -1)
			if !strings.Contains(ouput, expected) {
				t.Fatalf("unexpected output: %v", string(out))
			}
		})
	}
}

func TestGenerateAction(t *testing.T) {
	testCases := []struct {
		name       string
		actionName string
		output     string
	}{
		{
			name:       "No action name",
			actionName: "",
			output:     "Usage: generate action <folder/action>\n",
		},
		{
			name:       "r",
			actionName: "some/action",
			output:     "Action files created successfully✅\n",
		},
		{
			name:       "withouth slash",
			actionName: "action",
			output:     "Usage: generate action <folder/action>\n",
		},
		{
			name:       "two slashes",
			actionName: "some/other/action",
			output:     "Action files created successfully✅\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".", "gen", "action", tc.actionName)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running command: %v", err)
			}

			//Remove the files created in the tests
			if err := os.RemoveAll("internal/some"); err != nil {
				t.Fatalf("error removing files: %v", err)
			}

			if string(out) != tc.output {
				t.Fatalf("unexpected output: %v got: %v", string(out), tc.output)
			}
		})
	}
}
