package server_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/leapkit/leapkit/core/server"
)

func TestLogger(t *testing.T) {
	t.Run("base logger", func(t *testing.T) {
		output := bytes.NewBufferString("")
		current := os.Stderr

		r, testSrdErr, _ := os.Pipe()
		os.Stderr = testSrdErr

		t.Cleanup(func() {
			output.Reset()
			os.Stderr = current
		})

		s := server.New()
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/sssss", nil)

		s.Handler().ServeHTTP(resp, req)

		testSrdErr.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		if !strings.Contains(buf.String(), "url=/sssss") {
			t.Fatal("expected logs to contain plain url=/sssss")
		}
	})
}
