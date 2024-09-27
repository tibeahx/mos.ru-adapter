package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tibeahx/mos.ru-adapter/internal/config"
	"github.com/tibeahx/mos.ru-adapter/internal/handler"

	"go.uber.org/zap"
)

type Server struct {
	logger     *zap.SugaredLogger
	listenAddr string
	httpServer *http.Server
}

func NewServer(
	cfg *config.Config,
	handler *handler.Handler,
	logger *zap.SugaredLogger,
) *Server {
	httpServer := &http.Server{
		Handler:           handler.Router,
		Addr:              cfg.SrvListenAddr,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    4096,
	}

	return &Server{
		listenAddr: cfg.SrvListenAddr,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Run() error {
	if s == nil {
		return fmt.Errorf("server is nil")
	}
	if s.httpServer == nil {
		return fmt.Errorf("http server is nil")
	}
	s.logger.Infof("listening on %s:", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Infof("shutting down...")
	return s.httpServer.Shutdown(ctx)
}
