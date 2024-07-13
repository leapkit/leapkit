package sqlite

import (
	"fmt"
	"net/url"
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
	u, err := url.Parse(a.url)
	if err != nil {
		return fmt.Errorf("error parsing database URL: %w", err)
	}

	_, err = os.Create(u.Path)
	if err != nil {
		return fmt.Errorf("error creating database file %s: %w", u.Path, err)
	}

	return nil
}

// Drop sqlite database by removing the database file.
func (a *manager) Drop() error {
	u, err := url.Parse(a.url)
	if err != nil {
		return fmt.Errorf("error parsing database URL: %w", err)
	}

	err = os.Remove(u.Path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error dropping database: %w", err)
	}

	return nil
}
