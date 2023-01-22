package api

import (
	"go-bank-api/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  sqlc.Store
	router *gin.Engine
}

func NewServer(store sqlc.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	api := router.Group("/api")

	// register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// apply validCurrency validator on currency field
		v.RegisterValidation("currency", validCurrency)
	}

	// accounts
	api.POST("/accounts", server.createAccount)
	api.GET("/accounts/:id", server.getAccountById)
	api.GET("/accounts", server.listAccounts)

	// transfer
	api.POST("/transfer", server.createTransfer)

	server.router = router
	return server
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
