package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Option for the session middleware
type Option func(*sessions.CookieStore)

// Set the domain for the application session
// This is useful when you want to share the session
// between subdomains.
func WithDomain(domain string) Option {
	return func(store *sessions.CookieStore) {
		store.Options.Domain = domain
	}
}

// WithSecure value for the Secure flag on the session cookie.
func WithSecure(secure bool) Option {
	return func(store *sessions.CookieStore) {
		store.Options.Secure = secure
	}
}

// WithSameSite value for the SameSite option on the session cookie.
func WithSameSite(sameSite http.SameSite) Option {
	return func(store *sessions.CookieStore) {
		store.Options.SameSite = sameSite
	}
}
