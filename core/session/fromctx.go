package session

import (
	"context"

	"github.com/gorilla/sessions"
)

// FromCtx returns the session from the context.
func FromCtx(ctx context.Context) *sessions.Session {
	return ctx.Value(ctxKey).(*sessions.Session)
}
