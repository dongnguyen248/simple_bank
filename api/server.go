package api

import (
	db "github.com/dongnguyen248/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
