package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/log"
	"github.com/jljl1337/xpense/internal/server"
)

func main() {
	env.MustSetConstants()

	log.SetCustomLogger()

	// Start the server with graceful shutdown
	server, err := server.NewServer()
	if err != nil {
		slog.Error("Failed to create server: " + err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error { return server.Start() })
	g.Go(func() error {
		<-gCtx.Done()
		return server.Stop(context.Background())
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server error: " + err.Error())
	}
}
