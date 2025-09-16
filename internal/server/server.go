package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/repository"
	"github.com/jljl1337/xpense/internal/server/handler"
	"github.com/jljl1337/xpense/internal/server/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type Server struct {
	db         *sql.DB
	queries    *repository.Queries
	httpServer *http.Server
}

func NewServer(db *sql.DB) *Server {
	queries := repository.New(db)
	mux := http.NewServeMux()

	apiMux := http.NewServeMux()

	authService := service.NewAuthService(queries)
	userService := service.NewUserService(queries)
	bookService := service.NewBookService(queries)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	bookHandler := handler.NewBookHandler(bookService)

	healthHandler.RegisterRoutes(apiMux)
	authHandler.RegisterRoutes(apiMux)
	userHandler.RegisterRoutes(apiMux)
	bookHandler.RegisterRoutes(apiMux)

	middlewareProvider := middleware.NewMiddlewareProvider(authService)

	stack := middleware.CreateStack(
		middlewareProvider.CORS(),
		middlewareProvider.Logging(),
		middlewareProvider.Auth(),
	)

	mux.Handle("/api/", http.StripPrefix("/api", stack(apiMux)))

	// Serve static site
	webHandler := handler.NewWebHandler()
	mux.HandleFunc("/", webHandler.ServeSite)

	return &Server{
		db:      db,
		queries: queries,
		httpServer: &http.Server{
			Addr:    ":" + env.Port,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	slog.Info("Starting server on port " + env.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stopping server")

	if err := s.db.Close(); err != nil {
		slog.Error("Failed to close database connection: " + err.Error())
	}

	return s.httpServer.Shutdown(ctx)
}
