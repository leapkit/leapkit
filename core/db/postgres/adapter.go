package postgres

import (
	"database/sql"
	"regexp"
)

var (
	// urlExp is the regular expression to extract the database name
	// and the user credentials from the database URL.
	urlExp = regexp.MustCompile(`postgres:\/\/([^:]+):([^@]+)@([^:]+):(\d+)\/([^?]+).*`)
)

// adapter for the sqlite database it includes the connection
// to perform the framework operations.
type adapter struct {
	conn *sql.DB
}

// New sqlite adapter with the passed connection.
func New(conn *sql.DB) *adapter {
	return &adapter{
		conn: conn,
	}
}
