package main

import (
	"database/sql"
	"log"

	"github.com/dongnguyen248/simple_bank/api"
	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/dongnguyen248/simple_bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbDriver := config.DBDriver
	dbSource := config.DBSource
	serverAddress := config.ServerAddress
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
