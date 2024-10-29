package server

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"io/fs"

	"github.com/leapkit/leapkit/core/server/internal/response"
)

// Router is the interface that wraps the basic methods for a router
type Router interface {
	// Use allows to specify a middleware that should be executed for all the handlers
	Use(middleware ...Middleware)

	// ResetMiddleware clears the list of middleware on the router by setting the baseMiddleware.
	ResetMiddleware()

	// Handle allows to register a new handler for a specific pattern
	Handle(pattern string, handler http.Handler)

	// HandleFunc allows to register a new handler function for a specific pattern
	HandleFunc(pattern string, handler http.HandlerFunc)

	// Folder allows to serve static files from a directory
	Folder(prefix string, fs fs.FS)

	// Group allows to create a new group of routes with a common prefix
	Group(prefix string, fn func(Router))
}

// router is a group of routes with a common prefix and middleware
// that should be executed for all the handlers in the group
type router struct {
	prefix     string
	mux        *http.ServeMux
	middleware []Middleware
	rootSet    bool
}

// Use allows to specify a middleware that should be executed for all the handlers
// in the group
func (rg *router) Use(middleware ...Middleware) {
	rg.middleware = append(rg.middleware, middleware...)
}

// ResetMiddleware clears the list of middleware on the router by setting the baseMiddleware.
func (rg *router) ResetMiddleware() {
	rg.middleware = baseMiddleware
}

// Handle allows to register a new handler for a specific pattern
// in the group with the middleware that should be executed for the handler
// specified in the group.
func (rg *router) Handle(pattern string, handler http.Handler) {
	method := ""
	route := pattern

	if parts := strings.Split(pattern, " "); len(parts) > 1 {
		method = parts[0]
		route = parts[1]
	}

	pattern = fmt.Sprintf("%s %s", method, path.Join(rg.prefix, route))
	pattern = strings.Trim(pattern, " ")

	// When this route is set we mark the rootSet as true
	rg.rootSet = rg.rootSet || (pattern == "/")

	handler = recoverer(handler)

	// Wrapping with the middleware
	for i := len(rg.middleware) - 1; i >= 0; i-- {
		handler = rg.middleware[i](handler)
	}

	rg.mux.Handle(pattern, handler)
}

// HandleFunc allows to register a new handler function for a specific pattern
// in the group with the middleware that should be executed for the handler
// specified in the group.
func (rg *router) HandleFunc(pattern string, handler http.HandlerFunc) {
	rg.Handle(pattern, http.HandlerFunc(handler))
}

// Folder allows to serve static files from a directory
func (rg *router) Folder(prefix string, fs fs.FS) {
	rg.mux.Handle(
		fmt.Sprintf("GET %s/", path.Join(rg.prefix, prefix)),
		http.StripPrefix(prefix, http.FileServerFS(fs)),
	)
}

// Group allows to create a new group of routes with a common prefix
// and middleware that should be executed for all the handlers in the group
func (rg *router) Group(prefix string, rfn func(rg Router)) {
	group := &router{
		prefix:     path.Join(rg.prefix, prefix),
		mux:        rg.mux,
		middleware: rg.middleware,
	}

	rfn(group)
}

func (rg *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = &response.Writer{ResponseWriter: w}
	rg.mux.ServeHTTP(w, r)
}
