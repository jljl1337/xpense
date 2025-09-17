package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/jljl1337/xpense/internal/db"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/log"
	"github.com/jljl1337/xpense/internal/server"
)

func main() {
	env.SetConstants()

	log.SetCustomLogger()

	// Connect to the database
	dbInstance, err := db.NewDB(env.DbPath)
	if err != nil {
		slog.Error("Failed to connect to database: " + err.Error())
		return
	}

	// Migrate the database
	if err := db.Migrate(dbInstance); err != nil {
		slog.Error("Failed to migrate database: " + err.Error())
		return
	}

	// Start the server with graceful shutdown
	server := server.NewServer(dbInstance)

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
