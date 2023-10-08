package api

import (
	db "github.com/beatzoid/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves http requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new http server and sets up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on the specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
