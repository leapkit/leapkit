package generate_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/leapkit/leapkit/kit/generate"
)

func TestGenerateHandler(t *testing.T) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	err := os.Chdir(t.TempDir())
	if err != nil {
		t.Fatalf("error changing directory: %v", err)
	}

	err = generate.Handler("someHandler")
	if err != nil {
		t.Fatalf("error creating handler: %v", err)
	}

	// check the content of the file
	bb, err := os.ReadFile("internal/somehandler.go")
	if err != nil {
		t.Fatalf("handler was not created: %v", err)
	}

	contents := []string{
		"package internal",
		"func SomeHandler(",
	}

	for _, expected := range contents {
		if !bytes.Contains(bb, []byte(expected)) {
			t.Fatalf("content %v function not found", expected)
		}
	}
}
