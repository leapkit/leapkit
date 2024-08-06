package session_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/plush/v5"
	"github.com/leapkit/leapkit/core/server"
	"github.com/leapkit/leapkit/core/server/session"
)

func Test_Session_Setup(t *testing.T) {
	requestContext := context.Background()

	s := server.New(
		server.WithSession("session_test", "test"),
	)

	s.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(w, req)

			// capturing http.Request context
			requestContext = context.WithoutCancel(req.Context())
		})
	})

	s.Group("/values", func(rg server.Router) {
		rg.HandleFunc("GET /add/{value}/{$}", func(w http.ResponseWriter, r *http.Request) {
			sw := session.FromCtx(r.Context())
			sw.Values["value"] = r.PathValue("value")

			w.Write([]byte("OK"))
		})

		rg.HandleFunc("GET /clear/{$}", func(w http.ResponseWriter, r *http.Request) {
			sw := session.FromCtx(r.Context())
			for k := range sw.Values {
				delete(sw.Values, k)
			}

			w.Write([]byte("OK"))
		})

		rg.HandleFunc("GET /all/{$}", func(w http.ResponseWriter, r *http.Request) {
			sw := session.FromCtx(r.Context())
			v, _ := sw.Values["value"].(string)

			w.Write([]byte(v))
		})
	})

	s.Group("/flashes", func(rg server.Router) {
		rg.HandleFunc("GET /add/{value}/{$}", func(w http.ResponseWriter, r *http.Request) {
			sw := session.FromCtx(r.Context())
			sw.AddFlash(r.PathValue("value"))

			w.Write([]byte("OK"))
		})

		rg.HandleFunc("GET /all/{$}", func(w http.ResponseWriter, r *http.Request) {
			sw := session.FromCtx(r.Context())
			v := fmt.Sprint(sw.Flashes())

			w.Write([]byte(v))
		})

		rg.HandleFunc("GET /render/{$}", func(w http.ResponseWriter, r *http.Request) {
			valuer := r.Context().Value("valuer").(interface{ Values() map[string]any })
			result, _ := plush.Render(`<%= flash("_flash") %>`, plush.NewContextWith(valuer.Values()))

			w.Write([]byte(result))
		})
	})

	t.Run("session values", func(t *testing.T) {
		res := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/values/add/hello/", nil)
		s.Handler().ServeHTTP(res, req)

		if res.Body.String() != "OK" {
			t.Errorf("Expected 'OK', got '%s'", res.Body.String())
			return
		}

		sw := session.FromCtx(requestContext)
		if sw.Values["value"] != "hello" {
			t.Errorf("Expected 'hello', got '%s'", sw.Values["value"])
		}

		// Attempting 10 times to check if session values persist
		for range 10 {
			req = httptest.NewRequest(http.MethodGet, "/values/all/", nil)
			req = req.WithContext(requestContext)
			res.Body.Reset()

			s.Handler().ServeHTTP(res, req)

			if res.Body.String() != "hello" {
				t.Errorf("Expected 'OK', got '%s'", res.Body.String())
				return
			}
		}

		req = httptest.NewRequest(http.MethodGet, "/values/add/bar/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() != "OK" {
			t.Errorf("Expected 'OK', got '%s'", res.Body.String())
			return
		}

		req = httptest.NewRequest(http.MethodGet, "/values/clear/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() != "OK" {
			t.Errorf("Expected 'OK', got '%s'", res.Body.String())
			return
		}

		sw = session.FromCtx(requestContext)
		if len(sw.Values) > 0 {
			t.Errorf("Expected empty values, got '%v'", sw.Values)
			return
		}

		req = httptest.NewRequest(http.MethodGet, "/values/all/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() == "hello" {
			t.Errorf("Expected empty values, got '%s'", res.Body.String())
			return
		}

		requestContext = context.Background()
	})

	t.Run("session flashes", func(t *testing.T) {
		res := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/flashes/add/hello/", nil)
		s.Handler().ServeHTTP(res, req)

		responseHeader := res.Header()
		cookies, ok := responseHeader["Set-Cookie"]
		if !ok || len(cookies) != 1 {
			t.Fatal("No cookies. Header:", responseHeader)
		}

		if res.Body.String() != "OK" {
			t.Fatalf("Expected 'OK', got '%s'", res.Body.String())
		}

		sw := session.FromCtx(requestContext)
		if flashes, ok := sw.Values["_flash"].([]any); !ok {
			t.Fatalf("Expected non-empty flashes, got '%v'", flashes...)
		}

		req = httptest.NewRequest(http.MethodGet, "/flashes/all/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() != "[hello]" {
			t.Fatalf("Expected '[hello]', got '%s'", res.Body.String())
		}

		// Second attempt at the same endpoint to validate that there are no longer any flashes.
		req = httptest.NewRequest(http.MethodGet, "/flashes/all/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() == "hello" {
			t.Fatalf("Expected empty flashes, got '%s'", res.Body.String())
		}
	})

	t.Run("session flash helpers", func(t *testing.T) {
		res := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/flashes/add/hello/", nil)
		s.Handler().ServeHTTP(res, req)

		req = httptest.NewRequest(http.MethodGet, "/flashes/render/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() != "hello" {
			t.Fatalf("Expected 'hello' flashes, got '%s'", res.Body.String())
		}

		// Once the flash was called in the previous call, this should be removed from flashes.
		req = httptest.NewRequest(http.MethodGet, "/flashes/all/", nil)
		req = req.WithContext(requestContext)
		res.Body.Reset()

		s.Handler().ServeHTTP(res, req)

		if res.Body.String() != "[]" {
			t.Fatalf("Expected '[]', got '%s'", res.Body.String())
		}
	})
}
