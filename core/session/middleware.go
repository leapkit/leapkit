package session

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"
)

var ctxKey = "session"

func init() {
	// TODO: Look for a better place
	gob.Register(uuid.UUID{})
}

// InCtx puts the session in the context
var InCtx = Middleware

// Middleware that injects the session into the request context
// and also takes care of saving the session when the response is written
// to the client by wrapping the response writer.
func Middleware(secret, name string, options ...Option) func(http.Handler) http.Handler {
	store := sessions.NewCookieStore([]byte(secret))

	// Run the options on the store
	for _, option := range options {
		option(store)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, name)
			r = r.WithContext(context.WithValue(r.Context(), ctxKey, session))
			w = &saver{
				w:     w,
				req:   r,
				store: session,
			}

			type valueSetter interface {
				Set(key string, value interface{})
			}

			// Look for a valuer in the context and set the values for flash
			// and session so that they can be used in other components of the request.
			vlr, ok := r.Context().Value("valuer").(valueSetter)
			if ok {
				vlr.Set("flash", flashHelper(session))
				vlr.Set("session", func() *sessions.Session { return session })
			}

			if !ok {
				fmt.Println("no valuer in context")
			}

			next.ServeHTTP(w, r)
		})
	}
}
