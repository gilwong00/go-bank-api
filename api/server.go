package api

import (
	"go-bank-api/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *sqlc.Store
	router *gin.Engine
}

func NewServer(store *sqlc.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	api := router.Group("/api")

	// routes
	api.POST("/accounts", server.createAccount)
	api.GET("/accounts/:id", server.getAccountById)
	api.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
