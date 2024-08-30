package server

import (
	"net/http"
)

// defaultCatchAllHandler to log and return a 404 for all routes except the root route.
var defaultCatchAllHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		return
	}

	// TODO: Support custom 404 or have a default one.
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
})

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
	// if no catch-all or root route has been set
	// we use the default one
	if !s.rootSet {
		s.Handle("/", defaultCatchAllHandler)
	}

	return s.mux
}

func (s *mux) Addr() string {
	return s.host + ":" + s.port
}
