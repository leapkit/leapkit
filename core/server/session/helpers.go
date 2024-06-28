package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AddHelpers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := FromCtx(r.Context())

		// Add session helpers if there is a helperSetter in the context.
		rx, ok := r.Context().Value("renderer").(interface{ Set(string, any) })
		if ok {
			rx.Set("flash", flashHelper(session))
			rx.Set("session", func() *sessions.Session { return session })
		}

		next.ServeHTTP(w, r)
	})
}
