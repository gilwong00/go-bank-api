package grpcServer

import (
	"context"
	"fmt"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/token"
	"go-bank-api/pkg/util"
	"go-bank-api/rpc"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	// enable forward compability for the rpc calls before they are implemented
	rpc.UnimplementedBankServiceServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

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

func extractMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// HTTP request user agent
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
		// GRPC request user agent
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
		// HTTP request client IP
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			meta.ClientIP = clientIPs[0]
		}
	}
	// getting the client IP from the peer
	// this is for GRPC request
	if peer, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = peer.Addr.String()
	}
	return meta
}
