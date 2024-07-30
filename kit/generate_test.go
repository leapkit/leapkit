package main_test

import (
	"fmt"
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
		pckg       string
		funcName   string
	}{
		{
			name:       "No action name",
			actionName: "",
			output:     "Usage: generate action [action|folder/action]\n",
		},
		{
			name:       "one slash",
			actionName: "some/action",
			pckg:       "some",
			output:     "Action files created successfully✅\n",
			funcName:   "Action",
		},
		{
			name:       "withouth slash",
			actionName: "activity",
			pckg:       "internal",
			output:     "Action files created successfully✅\n",
			funcName:   "Activity",
		},
		{
			name:       "two slashes",
			actionName: "some/other/action",
			pckg:       "other",
			output:     "Action files created successfully✅\n",
			funcName:   "Action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".", "gen", "action", tc.actionName)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running command: %v", err)
			}

			if string(out) != tc.output {
				t.Fatalf("unexpected output: %v got: %v", string(out), tc.output)
			}

			if tc.actionName == "" {
				return
			}
			// Check if the file exists
			_, err = os.Stat("./internal/" + tc.actionName + ".go")
			if err != nil || os.IsNotExist(err) {
				t.Fatalf("file does not exist")
			}

			content, err := os.ReadFile("internal/" + tc.actionName + ".go")
			if err != nil {
				t.Fatalf("error reading file: %v", err)
			}

			if !strings.Contains(string(content), "package "+tc.pckg) {
				t.Fatalf("package name does not match")
			}

			if !strings.Contains(string(content), fmt.Sprintf(`rw.Render("%s.html")`, tc.actionName)) {
				t.Fatalf("html file does not match")
			}

			if !strings.Contains(string(content), fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request)", tc.funcName)) {
				t.Fatalf("function does not exist")
			}

			// Check if the file exists
			_, err = os.Stat("./internal/" + tc.actionName + ".html")
			if err != nil || os.IsNotExist(err) {
				t.Fatalf("file does not exist")
			}

			//Remove the files created in the tests
			if err := os.RemoveAll("internal/some"); err != nil {
				t.Fatalf("error removing files: %v", err)
			}

			if err := os.RemoveAll("internal/activity.go"); err != nil {
				t.Fatalf("error removing files: %v", err)
			}

			if err := os.RemoveAll("internal/activity.html"); err != nil {
				t.Fatalf("error removing files: %v", err)
			}
		})
	}
}

func TestGenerateHandler(t *testing.T) {
	testCases := []struct {
		name       string
		actionName string
		output     string
		pckg       string
		funcName   string
	}{
		{
			name:       "No action name",
			actionName: "",
			output:     "Usage: generate handler [name|folder/name]\n",
		},
		{
			name:       "one slash",
			actionName: "some/action",
			pckg:       "some",
			output:     "Handler file created successfully✅\n",
			funcName:   "Action",
		},
		{
			name:       "withouth slash",
			actionName: "activity",
			pckg:       "internal",
			output:     "Handler file created successfully✅\n",
			funcName:   "Activity",
		},
		{
			name:       "two slashes",
			actionName: "some/other/action",
			pckg:       "other",
			output:     "Handler file created successfully✅\n",
			funcName:   "Action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".", "gen", "handler", tc.actionName)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("error running command: %v", err)
			}

			if string(out) != tc.output {
				t.Fatalf("unexpected output: %v got: %v", string(out), tc.output)
			}

			if tc.actionName == "" {
				return
			}
			// Check if the file exists
			_, err = os.Stat("./internal/" + tc.actionName + ".go")
			if err != nil || os.IsNotExist(err) {
				t.Fatalf("file does not exist")
			}

			content, err := os.ReadFile("internal/" + tc.actionName + ".go")
			if err != nil {
				t.Fatalf("error reading file: %v", err)
			}

			if !strings.Contains(string(content), "package "+tc.pckg) {
				t.Fatalf("package name does not match")
			}

			if !strings.Contains(string(content), fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request)", tc.funcName)) {
				t.Fatalf("function does not exist")
			}

			//Remove the files created in the tests
			if err := os.RemoveAll("internal/some"); err != nil {
				t.Fatalf("error removing files: %v", err)
			}

			if err := os.RemoveAll("internal/activity.go"); err != nil {
				t.Fatalf("error removing files: %v", err)
			}
		})
	}
}
