package main

import (
	"Effective_Mobile/internal/config"
	"Effective_Mobile/internal/httpserver"
	"Effective_Mobile/internal/httpserver/routes"
	"Effective_Mobile/internal/logger"
	"Effective_Mobile/internal/storage/pg"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Effective Mobile API
// @version 1.0
// @description API для работы с данными о людях
// @host localhost:7007
// @BasePath /
func main() {
	cfg := config.MustLoad()
	logger.DebugEnabled = cfg.Debug

	logger.Info("Starting application")
	logger.Debug("Config loaded: %+v", cfg)

	storage, err := pg.New(&cfg.DsnPG)
	if err != nil {
		logger.Error("Failed to initialize storage: %v", err)
		os.Exit(1)
	}
	logger.Info("Storage initialized")

	router := routes.New(storage)
	server := httpserver.New(cfg.HTTPServer, *router)
	logger.Info("HTTP server initialized")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting HTTP server...")
		if err = server.Start(); err != nil {
			logger.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	logger.Info("Server started and listening")

	<-done
	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Stopping HTTP server")
	if err = server.Shutdown(ctx); err != nil {
		logger.Error("Failed to gracefully shutdown server: %v", err)
	}

	if err = storage.Close(); err != nil {
		logger.Error("Failed to close storage: %v", err)
	}

	logger.Info("Application stopped")
}
