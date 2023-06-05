package grpcServer

import (
	"fmt"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/token"
	"go-bank-api/pkg/util"
	"go-bank-api/rpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	// enable forward compability for the rpc calls before they are implemented
	rpc.UnimplementedBankServiceServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*GrpcServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &GrpcServer{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}

func dbUserToProtoUser(user db.User) *rpc.User {
	return &rpc.User{
		Username:          user.Username,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.LastName,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
