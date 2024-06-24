package db

import (
	"database/sql"
	"sync"
)

var (
	conn *sql.DB
	cmux sync.Mutex

	//DriverName defaults to postgres
	driverName = "postgres"
)

// ConnFn is the database connection builder function that
// will be used by the application based on the driver and
// connection string.
type ConnFn func() (*sql.DB, error)

// connectionOptions for the database
type connectionOption func()

// ConnectionFn is the database connection builder function that
// will be used by the application based on the driver and
// connection string. It opens the connection only once
// and return the same connection on subsequent calls.
func ConnectionFn(url string, opts ...connectionOption) ConnFn {
	return func() (cx *sql.DB, err error) {
		cmux.Lock()
		defer cmux.Unlock()

		if conn != nil && conn.Ping() == nil {
			return conn, nil
		}

		// Apply options before connecting to the database.
		for _, v := range opts {
			v()
		}

		conn, err = sql.Open(driverName, url)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}
}

// WithDriver allows to specify the driver to use driver defaults to
// postgres.
func WithDriver(name string) connectionOption {
	return func() {
		driverName = name
	}
}
