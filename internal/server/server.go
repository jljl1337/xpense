package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jmoiron/sqlx"

	"github.com/jljl1337/xpense/internal/db"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/http/handler"
	"github.com/jljl1337/xpense/internal/http/middleware"
	"github.com/jljl1337/xpense/internal/repository"
	"github.com/jljl1337/xpense/internal/service"
)

type Server struct {
	db         *sqlx.DB
	queries    *repository.Queries
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
	queries := repository.New(dbInstance)
	mux := http.NewServeMux()

	apiMux := http.NewServeMux()

	authService := service.NewAuthService(queries)
	userService := service.NewUserService(queries)
	bookService := service.NewBookService(queries)
	categoryService := service.NewCategoryService(queries)
	paymentMethodService := service.NewPaymentMethodService(queries)
	expenseService := service.NewExpenseService(queries)

	healthHandler := handler.NewHealthHandler()
	versionHandler := handler.NewVersionHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	bookHandler := handler.NewBookHandler(bookService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodService)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	healthHandler.RegisterRoutes(apiMux)
	versionHandler.RegisterRoutes(apiMux)
	authHandler.RegisterRoutes(apiMux)
	userHandler.RegisterRoutes(apiMux)
	bookHandler.RegisterRoutes(apiMux)
	categoryHandler.RegisterRoutes(apiMux)
	paymentMethodHandler.RegisterRoutes(apiMux)
	expenseHandler.RegisterRoutes(apiMux)

	middlewareProvider := middleware.NewMiddlewareProvider(authService)

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
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	if env.BackupCronSchedule != "" && env.BackupDbPath != "" {
		_, err = scheduler.NewJob(
			gocron.CronJob(
				env.BackupCronSchedule,
				false,
			),
			gocron.NewTask(
				func() {
					slog.Info("Starting database backup")
					start := time.Now()
					if err := db.BackupToFile(dbInstance, env.BackupDbPath); err != nil {
						slog.Error("Failed to backup database: " + err.Error())
						return
					}
					slog.Info("Database backup completed in " + time.Since(start).String())
				},
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create cron job: %w", err)
		}
	} else {
		slog.Warn("Database backup cron job not scheduled")
	}

	return &Server{
		db:      dbInstance,
		queries: queries,
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
