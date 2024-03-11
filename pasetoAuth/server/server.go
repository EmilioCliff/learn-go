package server

import (
	"fmt"
	"pasetoAuth/token"

	"github.com/gin-gonic/gin"
)

type Server struct {
	tokenMaker *token.PasetoMaker
	router     *gin.Engine
}

func NewServer(address string) (*Server, error) {
	tokenMaker, err := token.NewPaseto("abcdefghijkl12345678901234567890")
	if err != nil {
		return nil, fmt.Errorf("Couldnt open tokenmaker %w", err)
	}
	server := &Server{
		tokenMaker: tokenMaker,
	}

	server.setRoutes()
	server.router.Run(address)
	return server, nil
}

func (server *Server) setRoutes() {
	router := gin.Default()

	auth := router.Group("/").Use(authMiddleware(*server.tokenMaker))
	router.POST("/login", server.login)
	auth.DELETE("/delete/:id", server.deleteUser)
	router.POST("/create", server.createUser)

	server.router = router
}
