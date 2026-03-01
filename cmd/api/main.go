package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-github/internal/server"
)

func main() {
	srv := server.New()

	// Start server in a goroutine
	go func() {
		if err := srv.Run("8080"); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	slog.Info("server started", "port", 8080)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Graceful shutdown with 30s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.GracefulShutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}
