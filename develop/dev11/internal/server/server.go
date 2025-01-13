package server

import (
	"fmt"
	"net/http"
	"time"

	"dev11/internal/database"
)

type Server struct {
	db database.Service
}

func NewServer(host string, port int) *http.Server {
	NewServer := &Server{
		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
