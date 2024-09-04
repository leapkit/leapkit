package generate_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/leapkit/leapkit/kit/internal/generate"
)

func TestGenerateHandler(t *testing.T) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	tcases := []struct {
		name    string
		input   string
		file    string
		expects []string
	}{
		{
			name:  "simple case",
			input: "someHandler",
			file:  "internal/somehandler.go",
			expects: []string{
				"package internal",
				"func SomeHandler(",
			},
		},

		{
			name:  "simple folder",
			input: "some/handler",
			file:  "internal/some/handler.go",
			expects: []string{
				"package some",
				"func Handler(",
			},
		},

		{
			name:  "simple folder with ext",
			input: "some/handler.go",
			file:  "internal/some/handler.go",
			expects: []string{
				"package some",
				"func Handler(",
			},
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.name, func(t *testing.T) {
			err := os.Chdir(t.TempDir())
			if err != nil {
				t.Fatalf("error changing directory: %v", err)
			}

			err = generate.Handler(tcase.input)
			if err != nil {
				t.Fatalf("error creating handler: %v", err)
			}

			// check the content of the file
			bb, err := os.ReadFile(tcase.file)
			if err != nil {
				t.Fatalf("handler was not created: %v", err)
			}

			for _, expected := range tcase.expects {
				if !bytes.Contains(bb, []byte(expected)) {
					t.Log(string(bb))
					t.Fatalf("content %v function not found", expected)
				}
			}
		})
	}
}
