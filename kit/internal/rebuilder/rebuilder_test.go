package rebuilder_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/leapkit/leapkit/kit/internal/rebuilder"
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

	defer os.Remove("Procfile")

	t.Run("Correct", func(t *testing.T) {
		procfile()

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(time.Second)
			cancel()
		}()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Serve() returned an error: %v", err)
		}
	})

	t.Run("Correct - Skip invalid commands in Procfile", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := os.WriteFile("Procfile", []byte("web \napp : echo 'Hello, world!'"), 0644)
		if err != nil {
			t.Fatalf("Failed to write Procfile: %v", err)
		}

		go func() {
			time.Sleep(time.Second)
			cancel()
		}()

		if err := rebuilder.Serve(ctx); err != nil {
			t.Errorf("Expected nil, got '%v'", err)
		}
	})

	t.Run("Correct - Skip additional duplicated commands", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := os.WriteFile("Procfile", []byte("web: echo 'Hello, world!'\nweb: echo 'duplicated!'"), 0644)
		if err != nil {
			t.Fatalf("Failed to write Procfile: %v", err)
		}

		r, w, _ := os.Pipe()

		current := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = current
		})

		go func() {
			time.Sleep(time.Second)
			cancel()
		}()

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

	t.Run("Correct - Reloading commands", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := os.WriteFile("file.go", []byte("package main"), 0644)
		if err != nil {
			t.Fatalf("Failed to write custom_file.go: %v", err)
		}

		r, w, _ := os.Pipe()

		current := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = current
		})

		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Remove("file.go")
		}()

		go func() {
			time.Sleep(time.Second)
			cancel()
		}()

		fmt.Println("running the app")
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
