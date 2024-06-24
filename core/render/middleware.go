package render

import (
	"context"
	"io/fs"
	"net/http"
)

// InCtx puts the render engine in the context
// so the handlers can use it, it also sets a few
// other values that are useful for the handlers.
var InCtx = Middleware

// Middleware puts the render engine in the context
// so the handlers can use it, it also sets a few
// other values that are useful for the handlers.
func Middleware(templates fs.FS, options ...Option) func(http.Handler) http.Handler {
	engine := NewEngine(templates, options...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "renderer", engine.HTML(w))
			ctx = context.WithValue(ctx, "renderEngine", engine)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
