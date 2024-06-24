package db

import (
	"strings"

	"github.com/leapkit/core/db/postgres"
	"github.com/leapkit/core/db/sqlite"
)

// Manager is the interface that wraps the basic methods to
// create and drop a database.
type Manager interface {
	Create() error
	Drop() error
}

// Create a new database based on the passed URL.
func Create(url string) error {
	var adapter Manager = sqlite.NewManager(url)
	if strings.Contains(url, "postgres") {
		adapter = postgres.NewManager(url)
	}

	return adapter.Create()
}

// Drop a database based on the passed URL.
func Drop(url string) error {
	var adapter Manager = sqlite.NewManager(url)
	if strings.Contains(url, "postgres") {
		adapter = postgres.NewManager(url)
	}

	return adapter.Drop()
}
