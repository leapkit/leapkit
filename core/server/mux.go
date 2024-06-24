package server

import (
	"net/http"
)

// Rood routeGroup is a group of routes with a common prefix and middleware
// it also has a host and port as well as a Start method as it is the root of the server
// that should be executed for all the handlers in the group.
type mux struct {
	*router

	host string
	port string
}

// New creates a new server with the given options and default middleware.
func New(options ...Option) *mux {
	ss := &mux{
		router: &router{
			prefix:     "",
			mux:        http.NewServeMux(),
			middleware: baseMiddleware,
		},

		host: "0.0.0.0",
		port: "3000",
	}

	for _, option := range options {
		option(ss)
	}

	return ss
}

func (s *mux) Router() Router {
	return s.router
}

func (s *mux) Handler() http.Handler {
	return s.mux
}

func (s *mux) Addr() string {
	return s.host + ":" + s.port
}
