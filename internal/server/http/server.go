package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/VadimGossip/crudFinManager/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(netConfig config.NetServerConfig, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           netConfig.Host + ":" + strconv.Itoa(netConfig.Port),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
