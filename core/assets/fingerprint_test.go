package assets_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/leapkit/leapkit/core/assets"
)

func TestFingerprint(t *testing.T) {
	assetsPath := "/files/"
	m := assets.NewManager(fstest.MapFS{
		"main.js":        {Data: []byte("AAA")},
		"other/main.js":  {Data: []byte("AAA")},
		"custom/main.go": {Data: []byte("AAAA")},
	}, assetsPath)

	t.Run("is deterministic", func(t *testing.T) {
		a, _ := m.PathFor("files/main.js")
		b, _ := m.PathFor("files/main.js")
		if a != b {
			t.Errorf("Expected %s to equal %s", a, b)
		}

		if !strings.Contains(a, assetsPath) {
			t.Errorf("Expected %s to have %s prefix", a, assetsPath)
		}

		a, _ = m.PathFor("files/other/main.js")
		b, _ = m.PathFor("files/other/main.js")
		if a != b {
			t.Errorf("Expected %s to equal %s", a, b)
		}

		if !strings.Contains(a, assetsPath) {
			t.Errorf("Expected %s to have %s prefix", a, assetsPath)
		}
	})

	t.Run("adds starting slash", func(t *testing.T) {
		a, err := m.PathFor("files/main.js")
		if err != nil {
			t.Fatal(err)
		}

		b, err := m.PathFor("/files/main.js")
		if err != nil {
			t.Fatal(err)
		}

		if a != b {
			t.Errorf("Expected %s to equal %s", a, b)
		}
	})

	t.Run("adds starting /files", func(t *testing.T) {
		a, _ := m.PathFor("main.js")
		t.Log(a)
		if !strings.HasPrefix(a, "/files") {
			t.Errorf("Expected %s to start with /files", a)
		}
	})

	t.Run("respects folders", func(t *testing.T) {
		a, err := m.PathFor("main.js")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.HasPrefix(a, "/files/main-") {
			t.Errorf("Expected %s to contain /files/other/main-<hash>", a)
		}

		b, _ := m.PathFor("files/other/main.js")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.HasPrefix(b, "/files/other/main-") {
			t.Errorf("Expected %s to contain /files/other/main-<hash>", b)
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

	t.Run("file does not exist", func(t *testing.T) {
		a, err := m.PathFor("custom/main.go")
		if err == nil {
			t.Errorf("File must not exists: %s", a)
		}
	})

	t.Run("manager looks for files in root", func(t *testing.T) {
		m := assets.NewManager(fstest.MapFS{
			"main.js":       {Data: []byte("AAA")},
			"other/main.js": {Data: []byte("AAA")},
		}, "")

		a, err := m.PathFor("main.js")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.HasPrefix(a, "/") {
			t.Errorf("Expected %s to start with /", a)
		}

		b, err := m.PathFor("other/main.js")
		if err != nil {
			t.Fatal(err)
		}

		if !strings.HasPrefix(b, "/") {
			t.Errorf("Expected %s to start with /", b)
		}
	})

	t.Run("PathFor normalize file path", func(t *testing.T) {
		cases := []struct {
			servingPath string
			name        string
			expected    string
		}{
			{"public", "main.js", "/public/main-"},
			{"public", "/main.js", "/public/main-"},
			{"public", "public/main.js", "/public/main-"},
			{"public", "/public/main.js", "/public/main-"},
			{"public", "/other/main.js", "/public/other/main-"},
			{"public", "other/main.js", "/public/other/main-"},

			{"/public", "main.js", "/public/main-"},
			{"/public", "/main.js", "/public/main-"},
			{"/public", "public/main.js", "/public/main-"},
			{"/public", "/public/main.js", "/public/main-"},
			{"/public", "public/other/main.js", "/public/other/main-"},
			{"/public", "/public/other/main.js", "/public/other/main-"},

			{"/public/", "main.js", "/public/main-"},
			{"/public/", "/main.js", "/public/main-"},
			{"/public/", "public/main.js", "/public/main-"},
			{"/public/", "/public/main.js", "/public/main-"},
			{"/public/", "public/other/main.js", "/public/other/main-"},
			{"/public/", "/public/other/main.js", "/public/other/main-"},

			{"public/other", "main.js", "/public/other/main-"},
			{"public/other", "/main.js", "/public/other/main-"},
			{"public/other", "public/other/main.js", "/public/other/main-"},
			{"public/other", "/public/other/main.js", "/public/other/main-"},
			{"/public/other", "public/other/main.js", "/public/other/main-"},
			{"/public/other", "/public/other/main.js", "/public/other/main-"},

			{"/public/other", "main.js", "/public/other/main-"},
			{"/public/other", "/main.js", "/public/other/main-"},
			{"/public/other", "public/other/main.js", "/public/other/main-"},
			{"/public/other", "/public/other/main.js", "/public/other/main-"},

			{"/public/other/", "main.js", "/public/other/main-"},
			{"/public/other/", "/main.js", "/public/other/main-"},
			{"/public/other/", "public/other/main.js", "/public/other/main-"},
			{"/public/other/", "/public/other/main.js", "/public/other/main-"},
		}

		for i, c := range cases {
			manager := assets.NewManager(fstest.MapFS{
				"main.js":       {Data: []byte("AAA")},
				"other/main.js": {Data: []byte("AAA")},
			}, c.servingPath)

			result, err := manager.PathFor(c.name)
			if err != nil {
				t.Errorf("%d, Expected no error, got %s", i, err)
			}

			if !strings.HasPrefix(result, c.expected) {
				t.Errorf("%d, Expected %s to start with /public/", i, result)
			}
		}
	})
}
