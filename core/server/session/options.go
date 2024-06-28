package session

import "github.com/gorilla/sessions"

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
