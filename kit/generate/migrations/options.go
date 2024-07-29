package migrations

import (
	"fmt"
)

// Option is a function that applies an option
type Option func() error

// Apply applies the options
func Apply(options ...Option) error {
	for _, option := range options {
		if err := option(); err != nil {
			return fmt.Errorf("error applying migration option: %w", err)
		}
	}

	return nil
}

// UseFolder sets the folder for migrations
func UseFolder(f string) Option {
	return func() error {
		folder = f

		return nil
	}
}
