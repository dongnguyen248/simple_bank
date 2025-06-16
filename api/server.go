package api

import (
	"fmt"

	"github.com/dongnguyen248/simple_bank/token"
	"github.com/dongnguyen248/simple_bank/util"

	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	Router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)
	router.GET("/users/:user_name", server.getUser)
	router.GET("/users", server.listUsers)
	router.PUT("/users/change_password", server.changePassword)
	router.DELETE("/users/:user_name", server.deleteUser)
	router.POST("/users/login", server.loginUser)
	server.Router = router
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
