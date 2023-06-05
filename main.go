package main

import (
	"context"
	"database/sql"
	"go-bank-api/api"
	"go-bank-api/grpcServer"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
	"go-bank-api/rpc"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGrpcGatewayServer(config, store)
	startGrpcServer(config, store)
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

func runGrpcGatewayServer(config util.Config, store db.Store) {
	server, err := grpcServer.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create grpc server:", err)
	}
	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				// this keeps fields the same as it is defined in the proto file
				// for now set it to false
				UseProtoNames: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	// cancelling a context prevents the system from doing unnecessary work
	defer cancel()
	err = rpc.RegisterBankServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register bank service handler", err)
	}
	mux := http.NewServeMux()
	/*
		mux will receive HTTP request from clients. In order to convert the request
		into gRPC format, we will reroute all the request to the grpcMux
	*/
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("gRPChttp gateway server starting at %s:", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start gRPC http gateway server", err)
	}
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
