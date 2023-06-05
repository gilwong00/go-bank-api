package api

import (
	"database/sql"
	"fmt"
	"go-bank-api/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiration time.Time `json:"accessTokenExpiration"`
}

func (s *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	//verify refresh token is valid
	tokenPayload, err := s.tokenMaker.ValidateToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	// fetch user
	session, err := s.store.GetSession(ctx, tokenPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	// session validation
	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	if session.Username != tokenPayload.Username {
		err := fmt.Errorf("invalid user session")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("incorrect session token")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("session has expired")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	newAccessToken, newAccessTokenPayload, err := s.tokenMaker.CreateToken(tokenPayload.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, renewAccessTokenResponse{
		AccessToken:           newAccessToken,
		AccessTokenExpiration: newAccessTokenPayload.ExpiredAt,
	})
}
