package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// TestQueries provides all database queries
var testQueries *Queries

// TestDB provides a connection to the test database
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Create a new query from the DB connection
	testQueries = New(testDB)

	// os.exit() will return the status code of the tests
	// to the OS (e.x. if the tests fail, exit code 1 will be returned)
	os.Exit(m.Run())
}
