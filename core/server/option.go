package server

import (
	"net/http"

	"github.com/leapkit/leapkit/core/session"
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
	return func(m *mux) {
		m.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session.New(secret, name, options...).Register(w, r)

				h.ServeHTTP(w, r)
			})
		})
	}
}
