package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-co-op/gocron/v2"
	"github.com/jmoiron/sqlx"

	"github.com/jljl1337/xpense/internal/cron"
	"github.com/jljl1337/xpense/internal/db"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/handler"
	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type Server struct {
	db         *sqlx.DB
	httpServer *http.Server
	scheduler  gocron.Scheduler
}

func NewServer() (*Server, error) {
	// Connect to the database
	dbInstance, err := db.NewDB(env.DbPath, env.DbBusyTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Migrate the database
	if err := db.Migrate(dbInstance); err != nil {
		dbInstance.Close()
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Serve the API
	mux := http.NewServeMux()

	apiMux := http.NewServeMux()

	endpointService := service.NewEndpointService(dbInstance)
	endpointHandler := handler.NewEndpointHandler(endpointService)
	endpointHandler.RegisterRoutes(apiMux)

	middlewareService := service.NewMiddlewareService(dbInstance)
	middlewareProvider := middleware.NewMiddlewareProvider(middlewareService)

	stack := middleware.CreateStack(
		middlewareProvider.CORS(),
		middlewareProvider.Logging(),
		middlewareProvider.Auth(),
	)

	mux.Handle("/api/", http.StripPrefix("/api", stack(apiMux)))

	// Serve the static site
	webHandler := handler.NewWebHandler()
	mux.HandleFunc("/", webHandler.ServeSite)

	// Create the scheduler
	scheduler, err := cron.NewScheduler(dbInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Server{
		db: dbInstance,
		httpServer: &http.Server{
			Addr:    ":" + env.Port,
			Handler: mux,
		},
		scheduler: scheduler,
	}, nil
}

func (s *Server) Start() error {
	slog.Info("Starting scheduler")
	s.scheduler.Start()
	slog.Info("Starting server on port " + env.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stopping server")

	slog.Info("Stopping scheduler")
	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop scheduler: %w", err)
	}

	slog.Info("Closing database connection")
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return s.httpServer.Shutdown(ctx)
}
