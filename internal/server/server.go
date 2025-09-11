package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/jljl1337/xpense/internal/repository"
	"github.com/jljl1337/xpense/internal/server/handler"
	"github.com/jljl1337/xpense/internal/server/middleware"
	"github.com/jljl1337/xpense/internal/service"
)

type Server struct {
	queries    *repository.Queries
	httpServer *http.Server
}

func NewServer(db *sql.DB) *Server {
	queries := repository.New(db)
	mux := http.NewServeMux()

	apiMux := http.NewServeMux()

	authService := service.NewAuthService(queries)
	userService := service.NewUserService(queries)

	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	healthHandler.RegisterRoutes(apiMux)
	authHandler.RegisterRoutes(apiMux)
	userHandler.RegisterRoutes(apiMux)

	middlewareProvider := middleware.NewMiddlewareProvider(authService)

	stack := middleware.CreateStack(
		middlewareProvider.CORS(),
		middlewareProvider.Logging(),
		middlewareProvider.Auth(),
	)

	mux.Handle("/api/", http.StripPrefix("/api", stack(apiMux)))

	// Serve static files (React app)
	mux.HandleFunc("/", handler.WebHandler)

	return &Server{
		queries: queries,
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	slog.Info("Starting server")
	return s.httpServer.ListenAndServe()
}
