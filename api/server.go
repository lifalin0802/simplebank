package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

// Server serves HTTP requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup service
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)

	//remove the last / to verify only accept params from body, that's page_id, page_size
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
