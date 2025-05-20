package test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// create a connection to the database
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = db.New(testDB)

	// run the tests
	code := m.Run()

	// close the database connection
	if err := testDB.Close(); err != nil {
		log.Fatal("cannot close db:", err)
	}

	os.Exit(code)
}
