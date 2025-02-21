package rebuilder_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/leapkit/leapkit/tools/dev/internal/rebuilder"
)

func TestServe(t *testing.T) {
	// Create a Procfile
	procfile := func() *os.File {
		if _, err := os.Stat("Procfile"); !os.IsNotExist(err) {
			os.Remove("Procfile")
		}

		f, err := os.Create("Procfile")
		if err != nil {
			t.Fatalf("Failed to create Procfile: %v", err)
		}

		f.WriteString("web : echo 'Hello, web!'\napp : echo 'Hello, app!'\n INVALID COMMAND\n foo : echo 'Hello, foo!'")

		return f
	}

	testFile := func() *os.File {
		if _, err := os.Stat("test"); !os.IsNotExist(err) {
			os.RemoveAll("test")
		}

		os.Mkdir("test", 0o755)
		f, err := os.Create("test/main.go")
		if err != nil {
			t.Fatalf("Failed to create test/main.go: %v", err)
		}

		content := "package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(\"Hello, world!\")\n}"
		f.WriteString(content)

		return f
	}

	defer os.Remove("Procfile")
	defer os.RemoveAll("test")

	t.Run("Correct", func(t *testing.T) {
		r, w, _ := os.Pipe()

		current := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = current
		})

		procfile()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Serve() returned an error: %v", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if !strings.Contains(buf.String(), "[kit] Starting app") {
			t.Errorf("Expected '[kit] Starting app' to be in the output, got '%v'", buf.String())
		}

		if !strings.Contains(buf.String(), "[kit] Shutting down...") {
			t.Errorf("Expected '[kit] Shutting down...' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Correct - Watching added and removed files or folders within sub path", func(t *testing.T) {
		stdout := os.Stdout
		stderr := os.Stderr

		r, w, _ := os.Pipe()
		os.Stdout = w
		os.Stderr = w

		t.Cleanup(func() {
			os.Stdout = stdout
			os.Stderr = stderr
		})

		procfile()
		err := os.WriteFile("Procfile", []byte("web : echo 'Hello, world!'"), 0o644)
		if err != nil {
			t.Fatalf("Failed to write Procfile: %v", err)
		}

		os.MkdirAll("test/sub/folder/path", 0o755)
		os.Create("test/sub/folder/path/main.go")

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		go func() {
			time.Sleep(20 * time.Millisecond)
			os.RemoveAll("test")
		}()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Serve() returned an error: %v", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if !strings.Contains(buf.String(), "[kit] Starting app") {
			t.Errorf("Expected '[kit] Starting app' to be in the output, got '%v'", buf.String())
		}

		if !strings.Contains(buf.String(), "Restarted...") {
			t.Errorf("Expected 'Restarted...' to be in the output, got '%v'", buf.String())
		}

		if !strings.Contains(buf.String(), "[kit] Shutting down...") {
			t.Errorf("Expected '[kit] Shutting down...' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Correct - Skip invalid commands in Procfile", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := os.WriteFile("Procfile", []byte("web \napp : echo 'Hello, world!'"), 0o644)
		if err != nil {
			t.Fatalf("Failed to write Procfile: %v", err)
		}

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Expected nil, got '%v'", err)
		}
	})

	t.Run("Correct - Skip additional duplicated commands", func(t *testing.T) {
		r, w, _ := os.Pipe()

		current := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = current
		})

		procfile()
		os.WriteFile("Procfile", []byte("web: echo 'Hello, world!'\nweb: echo 'duplicated!'"), 0o644)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Expected nil, got '%v'", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if strings.Contains(buf.String(), "duplicated!") {
			t.Errorf("Expected 'duplicated!' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Correct - Printing multiline logs", func(t *testing.T) {
		r, w, _ := os.Pipe()

		stdOut := os.Stdout
		stdErr := os.Stderr

		os.Stdout = w
		os.Stderr = w
		t.Cleanup(func() {
			os.Stdout = stdOut
			os.Stderr = stdErr
		})

		procfile()
		os.WriteFile("Procfile", []byte("app: go run test/main.go"), 0o644)

		testFile()
		content := "package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(`Multiline\nlogs`)\n}"
		os.WriteFile("test/main.go", []byte(content), 0o644)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Expected nil, got '%v'", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if !strings.Contains(buf.String(), "app |\033[0m Multiline\n") {
			t.Errorf("Expected 'app |\033[0m Multiline' to be in the output, got '%v'", buf.String())
		}

		if !strings.Contains(buf.String(), "app |\033[0m logs\n") {
			t.Errorf("Expected 'app |\033[0m logs' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Correct - Reloading commands", func(t *testing.T) {
		r, w, _ := os.Pipe()

		current := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = current
		})

		procfile()
		os.WriteFile("Procfile", []byte("test: go run test/main.go"), 0o644)

		testFile()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		go func() {
			time.Sleep(10 * time.Millisecond)
			content := "package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(\"Updated!\")\n}"
			os.WriteFile("test/main.go", []byte(content), 0o644)
		}()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Serve() returned an error: %v", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if !strings.Contains(buf.String(), "Restarted...") {
			t.Errorf("Expected 'Restarted...' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Correct - Debounce mechanism to avoid multiple restarts", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		procfile()
		os.WriteFile("Procfile", []byte("test: go run test/main.go"), 0o644)

		testFile()

		r, w, _ := os.Pipe()

		current := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = current
		})

		go func() {
			content := "package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(\"Updated!\")\n}"

			for range 5 {
				time.Sleep(10 * time.Millisecond)
				os.WriteFile("test/main.go", []byte(content), 0o644)
			}

			cancel()
		}()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Serve() returned an error: %v", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if strings.Count(buf.String(), "Restarted...") != 1 {
			t.Errorf("Expected 'Restarted...' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Correct - Command failed and should wait for refresh", func(t *testing.T) {
		procfile()
		os.WriteFile("Procfile", []byte("test: go run test/main.go"), 0o644)
		content := "package main\n\nimport ( \"fmt\"\n \"time\"\n\n)\n\nfunc main() {\n time.Sleep(10*time.Millisecond)\n panic(fmt.Sprint(\"Something went wrong!\"))\n}"
		os.WriteFile("test/main.go", []byte(content), 0o644)
		testFile()

		r, w, _ := os.Pipe()

		stdOut := os.Stdout
		stdErr := os.Stderr

		os.Stdout = w
		os.Stderr = w
		t.Cleanup(func() {
			os.Stdout = stdOut
			os.Stderr = stdErr
		})

		ctx, cancel := context.WithCancel(context.Background())

		// Refresh
		go func() {
			time.Sleep(100 * time.Millisecond)
			content = "package main\n\nimport ( \"fmt\"\n \"time\"\n\n)\n\nfunc main() {\n time.Sleep(10*time.Millisecond)\n}"
			os.WriteFile("test/main.go", []byte(content), 0o644)
		}()

		go func() {
			time.Sleep(300 * time.Millisecond)
			cancel()
		}()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Serve() returned an error: %v", err)
		}

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if strings.Count(buf.String(), "Restarted...") != 1 {
			t.Errorf("Expected 'Restarted...' to be in the output, got '%v'", buf.String())
		}
	})

	t.Run("Incorrect - Procfile not found", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		os.Remove("Procfile")

		go func() {
			time.Sleep(time.Second)
			cancel()
		}()

		if err := rebuilder.Serve(ctx); err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
