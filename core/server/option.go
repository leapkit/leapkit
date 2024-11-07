package server

import (
	"io/fs"
	"net/http"

	"github.com/leapkit/leapkit/core/assets"
	"github.com/leapkit/leapkit/core/server/session"
)

// Options for the server
type Option func(*mux)

// WithHost allows to specify the host to run the server at
// if not specified it defaults to 0.0.0.0
func WithHost(host string) Option {
	return func(s *mux) {
		s.host = host
	}
}

// WithPort allows to specify the port to run the server at
// when not specified it defaults to 3000
func WithPort(port string) Option {
	return func(s *mux) {
		s.port = port
	}
}

// WithSession allows to set the session within the application.
func WithSession(secret, name string, options ...session.Option) Option {
	sw := session.New(secret, name, options...)
	return func(m *mux) {
		m.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w, r = sw.Register(w, r)

				h.ServeHTTP(w, r)
			})
		})
	}
}

func WithAssets(embedded fs.FS) Option {
	manager := assets.NewManager(embedded)
	return func(m *mux) {
		m.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if vlr, ok := r.Context().Value("valuer").(interface{ Set(string, any) }); ok {
					vlr.Set("assetPath", manager.PathFor)
				}

				h.ServeHTTP(w, r)
			})
		})

		m.Folder(manager.HandlerPattern(), manager)
	}
}

func WithErrorMessage(status int, message string) Option {
	return func(m *mux) {
		errorMessageMap[status] = message
	}
}
