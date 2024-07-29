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

func WithPath(path string) Option {
	return func(store *sessions.CookieStore) {
		store.Options.Path = path
	}
}

// WithMaxAge sets the maximum age for the session cookie.
// A negative value means that the cookie is deleted when the browser is closed.
func WithMaxAge(maxAge int) Option {
	return func(store *sessions.CookieStore) {
		store.Options.MaxAge = maxAge
	}
}

// WithHTTPOnly sets the HttpOnly flag on the session cookie.
func WithHTTPOnly(httpOnly bool) Option {
	return func(store *sessions.CookieStore) {
		store.Options.HttpOnly = httpOnly
	}
}

// WithSecure sets the Secure flag on the session cookie.
func WithSecureFlag(secure bool) Option {
	return func(store *sessions.CookieStore) {
		store.Options.Secure = secure
	}
}
