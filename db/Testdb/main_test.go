package test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// create a connection to the database
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbDriver := config.DBDriver
	dbSource := config.DBSource
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
