package main

import (
	"database/sql"
	"go-bank-api/api"
	"go-bank-api/grpcServer"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
	"go-bank-api/rpc"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	// startGinServer(config, store)
	startGrpcServer(config, store)
}

func startGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.StartServer(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start http server:", err)
	}
}

func startGrpcServer(config util.Config, store db.Store) {
	server, err := grpcServer.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}
	gServer := grpc.NewServer()
	rpc.RegisterBankServiceServer(gServer, server)
	reflection.Register(gServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener:", err)
	}
	log.Printf("starting grpc server at %s:", listener.Addr().String())
	err = gServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start grpc server:", err)
	}
}
