package main

import (
	server "PVZ/internal/app"
	"PVZ/internal/app/closer"
	"PVZ/internal/config"
	"PVZ/internal/logger"
	"PVZ/internal/metrics"
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	servers, err := server.SetupServers(ctx, cfg)
	if err != nil {
		logger.Fatal("Failed to setup servers: %v", zap.Error(err))
	}

	metrics.Register()

	errChan := make(chan error, 1)
	go runHTTPServer(servers.HTTP, errChan)
	go runPrometheusServer(servers.Prometheus, errChan)
	closer.WaitForShutdown(ctx, errChan, servers)
}

func runHTTPServer(s *http.Server, errChan chan<- error) {
	logger.Info("Starting HTTP server on ", zap.String("address", s.Addr))
	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		errChan <- fmt.Errorf("HTTP server error: %w", err)
	}
}

func runPrometheusServer(s *http.Server, errChan chan<- error) {
	logger.Info("Starting Prometheus server", zap.String("address", s.Addr))
	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		errChan <- fmt.Errorf("prometheus server error: %w", err)
	}
}
