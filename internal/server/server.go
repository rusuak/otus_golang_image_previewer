package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	config Config
}

type Config struct {
	Addr              string
	ReadTimeout       int
	ReadHeaderTimeout int
}

func NewServer(config Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run(routers *mux.Router) error {
	server := &http.Server{
		Addr:              s.config.Addr,
		Handler:           routers,
		ReadTimeout:       time.Duration(s.config.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.config.ReadHeaderTimeout) * time.Second,
	}

	fmt.Println("Server is starting...")
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
