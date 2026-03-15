package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go-github/api" // Import generated docs
	internalmcp "go-github/internal/mcp"
	"go-github/internal/server"

	"golang.org/x/sync/errgroup"
)

// @title Home Lab API
// @version 1.0
// @description API for managing home lab devices and services
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a shared context that is cancelled on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := server.New()

	g, gctx := errgroup.WithContext(ctx)

	// Launch HTTP server goroutine.
	g.Go(func() error {
		slog.Info("http server started", "port", port)
		if err := srv.Run(port); err != nil {
			return err
		}
		return nil
	})

	// Launch MCP stdio server goroutine.
	g.Go(func() error {
		slog.Info("mcp server started", "transport", "stdio")
		return internalmcp.Run(gctx)
	})

	// Goroutine to trigger graceful HTTP shutdown when context is done.
	g.Go(func() error {
		<-gctx.Done()

		slog.Info("shutting down servers...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.GracefulShutdown(shutdownCtx); err != nil {
			slog.Error("http shutdown error", "error", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	slog.Info("all servers stopped gracefully")
}
