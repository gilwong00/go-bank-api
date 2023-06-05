package grpcServer

import (
	"context"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
	"go-bank-api/rpc"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GrpcServer) CreateUser(
	ctx context.Context,
	req *rpc.CreateUserRequest,
) (*rpc.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	payload := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		Email:          req.GetEmail(),
	}
	user, err := s.store.CreateUser(ctx, payload)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exist: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create error: %s", err)
	}
	return &rpc.CreateUserResponse{
		User: dbUserToProtoUser(user),
	}, nil
}
