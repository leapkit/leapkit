package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leapkit/core/server"
)

func TestRouter(t *testing.T) {

	s := server.New()

	s.Group("/", func(r server.Router) {
		r.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})

		r.Group("/api/", func(r server.Router) {
			r.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("This is the API!"))
			})

			r.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("API documentation!"))
			})

			r.Group("/v1/", func(r server.Router) {
				r.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("Welcome to the API v1!"))
				})

				r.Group("/users/", func(r server.Router) {
					r.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("Users list!"))
					})

					r.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
						w.Write([]byte("Hello users!"))
					})
				})
			})
		})
	})

	testCases := []struct {
		method string
		route  string
		body   string
		code   int
	}{
		{"GET", "/", "Hello, World!", http.StatusOK},
		{"GET", "/api/v1/users/hello", "Hello users!", http.StatusOK},
		{"GET", "/api/v1/users/", "Users list!", http.StatusOK},
		{"GET", "/api/v1/", "Welcome to the API v1!", http.StatusOK},
		{"GET", "/api/", "This is the API!", http.StatusOK},
		{"GET", "/api/docs", "API documentation!", http.StatusOK},
	}

	for _, tt := range testCases {
		t.Run(tt.route, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.route, nil)
			res := httptest.NewRecorder()
			s.Handler().ServeHTTP(res, req)

			if res.Code != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, res.Code)
			}

			if res.Body.String() != tt.body {
				t.Errorf("Expected body %s, got %s", tt.body, res.Body.String())
			}
		})
	}

}

func TestMiddleware(t *testing.T) {
	s := server.New()
	s.Use(server.InCtxMiddleware("customValue", "Hello, World!"))

	s.Group("/", func(r server.Router) {
		r.HandleFunc("GET /mw/{$}", func(w http.ResponseWriter, r *http.Request) {
			v := r.Context().Value("customValue").(string)
			w.Write([]byte(v))
		})

		r.Group("/without", func(r server.Router) {
			r.ResetMiddleware()

			r.HandleFunc("GET /mw/{$}", func(w http.ResponseWriter, r *http.Request) {
				v, ok := r.Context().Value("customValue").(string)
				if !ok {
					w.Write([]byte("customValue not found"))
					return
				}
				w.Write([]byte(v))
			})
		})

		r.Group("/other-with", func(r server.Router) {
			r.HandleFunc("GET /mw/{$}", func(w http.ResponseWriter, r *http.Request) {
				v := r.Context().Value("customValue").(string)
				w.Write([]byte(v + " (again)"))
			})
		})
	})

	testCases := []struct {
		description string
		pattern     string
		expected    string
	}{
		{"request to handler with middleware", "/mw/", "Hello, World!"},
		{"request to handler without middleware", "/without/mw/", "customValue not found"},
		{"request to other handler with middleware", "/other-with/mw/", "Hello, World! (again)"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tc.pattern, nil)
			res := httptest.NewRecorder()
			s.Handler().ServeHTTP(res, req)

			if res.Body.String() != tc.expected {
				t.Errorf("Expected body %s, got %s", tc.expected, res.Body.String())
			}
		})
	}
}
