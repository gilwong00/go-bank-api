package main

import (
	"database/sql"
	"go-bank-api/api"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start http server:", err)
	}
}
