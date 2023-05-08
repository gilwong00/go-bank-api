package api

import (
	"database/sql"
	"errors"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/token"
	"go-bank-api/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	// Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

/*
	server is a receiver
	you can declare a method with a server whos type is defined in the same package as the method
	we can also use pointer receivers like we do with server
	methods with pointer receivers can modify the value the receiver is pointing to
	The first is so that the method can modify the value that its receiver points to.
	The second is to avoid copying the value on each method call.
	This can be more efficient if the receiver is a large struct
*/

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet("auth_payload").(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountById(ctx *gin.Context) {
	var req getAccountByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	account, err := server.store.GetAccountById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet("auth_payload").(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=1,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload := ctx.MustGet("auth_payload").(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
