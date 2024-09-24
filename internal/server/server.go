package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-task/internal/config"
	"test-task/internal/handler"
	"time"

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
	s.logger.Infof("listening on %s:", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded {
				s.logger.Panic("shotdown timed out, forcing quit...")
			}
		}()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Fatal("failed to shutdown from ctx")
		}
	}()

	<-ctx.Done()

	return nil
}
