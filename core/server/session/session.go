package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/leapkit/leapkit/core/server/internal/writer"
)

// ctxKey is the value used to store the session
// into the http.Request context.
var ctxKey contextKey = "session"

// contextKey is the key type used to store the session
// into the http.Request context.
type contextKey string

func New(secret, name string, options ...Option) *session {
	store := sessions.NewCookieStore([]byte(secret))

	// Default options.
	store.Options.HttpOnly = true

	// TODO: Review these 2 options for production.
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteLaxMode

	// Run the options on the store
	for _, option := range options {
		option(store)
	}

	return &session{
		name:  name,
		store: store,
	}
}

type session struct {
	name  string
	store *sessions.CookieStore
}

// Register returns an *http.Request with the session set in its context and also
// a custom http.ResponseWriter implementation that will save the session after each HTTP call.
func (s *session) Register(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		fmt.Println(err, "session_name", s.name)
	}

	// Look for a valuer in the context and set the values for flash
	// and session so that they can be used in other components of the request.
	vlr, ok := r.Context().Value("valuer").(interface{ Set(string, any) })
	if ok {
		vlr.Set("flash", flashHelper(session))
		vlr.Set("session", func() *sessions.Session { return session })
	}

	r = r.WithContext(context.WithValue(r.Context(), ctxKey, session))

	w = &saver{
		ResponseWriter: &writer.ResponseWriter{ResponseWriter: w},
		req:            r,
		store:          session,
	}

	return w, r
}
