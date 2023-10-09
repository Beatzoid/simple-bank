package api

import (
	db "github.com/beatzoid/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves http requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new http server and sets up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// We don't use proxies so no need to trust any
	router.SetTrustedProxies(nil)

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/account", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)

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
