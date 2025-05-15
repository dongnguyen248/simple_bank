package test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/dongnguyen24/sqlc-demo/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *sqlc.Queries

func TestMain(m *testing.M) {
	// create a connection to the database
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)

	// run the tests
	code := m.Run()

	// close the database connection
	if err := conn.Close(); err != nil {
		log.Fatal("cannot close db:", err)
	}

	os.Exit(code)
}
