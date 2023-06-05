package grpcServer

import (
	"context"
	"database/sql"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/util"
	"go-bank-api/rpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GrpcServer) Login(
	ctx context.Context,
	req *rpc.LoginRequest,
) (*rpc.LoginResponse, error) {
	user, err := s.store.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}
	err = util.ValidatePassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invalid password: %s", err)
	}
	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}
	refreshToken, refreshTokenPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}
	return &rpc.LoginResponse{
		User:                   dbUserToProtoUser(user),
		SessionId:              session.ID.String(),
		AccessToken:            accessToken,
		AccessTokenExpiration:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}, nil
}
