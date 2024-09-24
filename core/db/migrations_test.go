package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/leapkit/leapkit/core/db"
	_ "github.com/mattn/go-sqlite3"
)

func TestSetup(t *testing.T) {
	td := t.TempDir()
	conn, err := sql.Open("sqlite", filepath.Join(td, "database.db"))
	if err != nil {
		t.Fatal(err)
	}

	adapter := db.NewMigrator(conn)
	err = adapter.Setup()
	if err != nil {
		t.Fatal(err)
	}

	var name string
	rows, err := conn.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='schema_migrations';")
	if err != nil {
		t.Fatal("schema_migrations table not found")
	}

	if !rows.Next() {
		t.Fatal("schema_migrations table not found")
	}

	err = rows.Scan(&name)
	if err != nil {
		t.Fatal(err)
	}

	if name != "schema_migrations" {
		t.Fatal("schema_migrations table not found")
	}
}

func TestRun(t *testing.T) {
	t.Run("migration not found", func(t *testing.T) {
		td := t.TempDir()
		conn, err := sql.Open("sqlite", filepath.Join(td, "database.db"))
		if err != nil {
			t.Fatal(err)
		}

		adapter := db.NewMigrator(conn)
		err = adapter.Setup()
		if err != nil {
			t.Fatal(err)
		}

		err = adapter.Run("20210101000000", "name", "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);")
		if err != nil {
			t.Fatal(err)
		}

		var name string
		row := conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users';")
		err = row.Scan(&name)
		if err != nil {
			t.Fatal("users table not found")
		}
	})

	t.Run("migration found", func(t *testing.T) {
		td := t.TempDir()
		conn, err := sql.Open("sqlite", filepath.Join(td, "database.db"))
		if err != nil {
			t.Fatal(err)
		}

		adapter := db.NewMigrator(conn)
		err = adapter.Setup()
		if err != nil {
			t.Fatal(err)
		}

		err = adapter.Run("20210101000000", "name", "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);")
		if err != nil {
			t.Fatal(err)
		}

		err = adapter.Run("20210101000000", "name", "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);")
		if err != nil {
			t.Fatal(err)
		}

		var name string
		row := conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users';")
		err = row.Scan(&name)
		if err != nil {
			t.Fatal("users table not found")
		}
	})
}
