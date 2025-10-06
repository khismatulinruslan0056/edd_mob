package httpserver

import (
	"Effective_Mobile/internal/config"
	"Effective_Mobile/internal/httpserver/routes"
	"Effective_Mobile/internal/logger"
	"context"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func New(cfg config.HTTPServer, router routes.Router) *HTTPServer {
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      &router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	logger.Info("httpserver: initialized with address %s (timeout: %s, idle timeout: %s)", cfg.Address, cfg.Timeout, cfg.IdleTimeout)

	return &HTTPServer{server: srv}
}

func (s *HTTPServer) Start() error {
	logger.Info("httpserver: starting server on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("httpserver: failed to start: %v", err)
		return err
	}
	logger.Info("httpserver: server stopped gracefully")
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	logger.Info("httpserver: shutting down server...")
	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Error("httpserver: shutdown error: %v", err)
	} else {
		logger.Info("httpserver: shutdown completed successfully")
	}
	return err
}
