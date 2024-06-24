package assets_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/leapkit/core/assets"
)

func TestFingerprint(t *testing.T) {
	m := assets.NewManager(fstest.MapFS{
		"main.js":       {Data: []byte("AAA")},
		"other/main.js": {Data: []byte("AAA")},
	})

	t.Run("is deterministic", func(t *testing.T) {
		a, _ := m.PathFor("public/main.js")
		b, _ := m.PathFor("public/main.js")
		if a != b {
			t.Errorf("Expected %s to equal %s", a, b)
		}

		if !strings.Contains(a, "/public/") {
			t.Errorf("Expected %s to have /public/ prefix", a)
		}

		a, _ = m.PathFor("public/other/main.js")
		b, _ = m.PathFor("public/other/main.js")
		if a != b {
			t.Errorf("Expected %s to equal %s", a, b)
		}

		if !strings.Contains(a, "/public/") {
			t.Errorf("Expected %s to have /public/ prefix", a)
		}
	})

	t.Run("adds starting slash", func(t *testing.T) {
		a, err := m.PathFor("public/main.js")
		if err != nil {
			t.Fatal(err)
		}

		b, err := m.PathFor("/public/main.js")
		if err != nil {
			t.Fatal(err)
		}

		if a != b {
			t.Errorf("Expected %s to equal %s", a, b)
		}
	})

	t.Run("adds starting /public", func(t *testing.T) {
		a, _ := m.PathFor("main.js")
		t.Log(a)
		if !strings.HasPrefix(a, "/public") {
			t.Errorf("Expected %s to start with /public", a)
		}
	})

	t.Run("respects folders", func(t *testing.T) {
		a, err := m.PathFor("main.js")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.HasPrefix(a, "/public/main-") {
			t.Errorf("Expected %s to contain /public/other/main-<hash>", a)
		}

		b, _ := m.PathFor("public/other/main.js")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.HasPrefix(b, "/public/other/main-") {
			t.Errorf("Expected %s to contain /public/other/main-<hash>", b)
		}

		if a == b {
			t.Errorf("Expected %s to not equal %s", a, b)
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		a, err := m.PathFor("foo.js")
		if err == nil {
			t.Errorf("File must not exists: %s", a)
		}
	})
}
