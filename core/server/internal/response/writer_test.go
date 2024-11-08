package response_test

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leapkit/leapkit/core/server"
)

// testHijackResponseWriter wraps the httptest.ResponseRecorder to implement Hijack interface.
type testHijackResponseWriter struct {
	*httptest.ResponseRecorder
}

func (t *testHijackResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func TestWriter(t *testing.T) {
	s := server.New(
		server.WithSession("test_secret", "test"),
	)

	s.HandleFunc("GET /flush/{$}", func(w http.ResponseWriter, _ *http.Request) {
		f, ok := w.(http.Flusher)
		if !ok {
			w.Write([]byte("Flush not supported"))
			return
		}

		f.Flush()

		w.Write([]byte("Flush supported!"))
	})

	s.HandleFunc("GET /hijack/{$}", func(w http.ResponseWriter, _ *http.Request) {
		h, ok := w.(http.Hijacker)
		if !ok {
			w.Write([]byte("Hijack not supported"))
			return
		}

		_, _, err := h.Hijack()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("Hijack supported!"))
	})

	resp := &testHijackResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}
	req := httptest.NewRequest(http.MethodGet, "/flush/", nil)
	s.Handler().ServeHTTP(resp, req)

	if resp.Body.String() == "Flush not supported" {
		t.Errorf("Expected 'Flush supported!', got 'Flush not supported'")
		return
	}

	resp = &testHijackResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}

	req = httptest.NewRequest(http.MethodGet, "/hijack/", nil)
	s.Handler().ServeHTTP(resp, req)

	if resp.Body.String() == "Hijack not supported" {
		t.Errorf("Expected 'Hijack supported!', got 'Hijack not supported'")
		return
	}
}
