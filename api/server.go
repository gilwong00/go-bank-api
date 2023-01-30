package api

import (
	"fmt"
	db "go-bank-api/pkg/db/sqlc"
	"go-bank-api/pkg/token"
	"go-bank-api/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// apply validCurrency validator on currency field
		v.RegisterValidation("currency", validCurrency)
	}

	server.initRoutes()
	return server, nil
}

func (server *Server) initRoutes() {
	router := gin.Default()
	api := router.Group("/api")

	// accounts
	api.POST("/accounts", server.createAccount)
	api.GET("/accounts/:id", server.getAccountById)
	api.GET("/accounts", server.listAccounts)

	// transfer
	api.POST("/transfer", server.createTransfer)

	//user
	api.POST("/user", server.createUser)
	api.POST("/user/auth", server.authUser)

	server.router = router
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
