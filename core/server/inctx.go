package server

import (
	"context"
	"net/http"
)

// InCtxMiddleware allows to specify a key/value that should be set on each
// request context. This is useful for services that could be used by the handlers.
func InCtxMiddleware(key string, value interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, value))

			next.ServeHTTP(w, r)
		})
	}
}
