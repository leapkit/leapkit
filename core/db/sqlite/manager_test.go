package sqlite_test

import (
	"os"
	"testing"

	"github.com/leapkit/leapkit/core/db/sqlite"
)

func TestCreate(t *testing.T) {
	f := t.TempDir()
	wd, _ := os.Getwd()

	os.Chdir(f)
	defer os.Chdir(wd)

	m := sqlite.NewManager("database.db?_fk=true")
	err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("database.db")
	if err != nil {
		t.Fatal("Did not create database.db", err)
	}

	_, err = os.Stat("database.db?_fk=true")
	if err == nil {
		t.Fatal("Created database.db?_fk=true", err)
	}
}

func TestDrop(t *testing.T) {
	f := t.TempDir()
	wd, _ := os.Getwd()

	os.Chdir(f)
	defer os.Chdir(wd)

	m := sqlite.NewManager("database.db?_fk=true")
	err := m.Create()
	if err != nil {
		t.Fatal(err)
	}

	err = m.Drop()
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("database.db")
	if err == nil {
		t.Fatal("Did not drop database.db", err)
	}
}
