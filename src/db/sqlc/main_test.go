package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/beatzoid/simple-bank/util"
	_ "github.com/lib/pq"
)

// TestQueries provides all database queries
var testQueries *Queries

// TestDB provides a connection to the test database
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Create a new query from the DB connection
	testQueries = New(testDB)

	// os.exit() will return the status code of the tests
	// to the OS (e.x. if the tests fail, exit code 1 will be returned)
	os.Exit(m.Run())
}
