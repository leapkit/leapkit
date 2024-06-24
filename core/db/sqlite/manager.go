package sqlite

import (
	"fmt"
	"os"
)

type manager struct {
	url string
}

func NewManager(url string) *manager {
	return &manager{
		url: url,
	}
}

// Create sqlite database file in the passed URL
func (a *manager) Create() error {
	_, err := os.Create(a.url)
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}

	return nil
}

// Drop sqlite database by removing the database file.
func (a *manager) Drop() error {
	err := os.Remove(a.url)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error dropping database: %w", err)
	}

	return nil
}
