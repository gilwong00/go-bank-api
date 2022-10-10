package main

import (
	"database/sql"
	"go-bank-api/api"
	"go-bank-api/sqlc"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://postgres:postgres@localhost:5432/bank_api?sslmode=disable"
	serverAddress = "0.0.0.0:5000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	store := sqlc.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(serverAddress)

	if err != nil {
		log.Fatal("Cannot start http server:", err)
	}
}
