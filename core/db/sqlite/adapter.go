package sqlite

import (
	"database/sql"
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
