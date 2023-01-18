package main

import (
	"database/sql"
	"go-bank-api/api"
	"go-bank-api/pkg/util"
	"go-bank-api/sqlc"
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

	store := sqlc.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start http server:", err)
	}
}
