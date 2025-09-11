package main

import (
	"log/slog"

	"github.com/jljl1337/xpense/internal/db"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/log"
	"github.com/jljl1337/xpense/internal/server"
)

func main() {
	env.LoadOptionalEnvFile()

	log.SetCustomLogger()

	// Migrate the database
	dbInstance, err := db.NewDB()
	if err != nil {
		slog.Error("Failed to connect to database: " + err.Error())
		return
	}

	if err := db.Migrate(dbInstance); err != nil {
		slog.Error("Failed to migrate database: " + err.Error())
		return
	}

	// queries := db.NewRepositoryQueries(dbInstance)
	// ctx := context.Background()
	// queries.CreateUser(ctx, repository.CreateUserParams{
	// 	ID:           ksuid.New().String(),
	// 	Email:        "user@example.com",
	// 	PasswordHash: "hashed_password",
	// 	CreatedAt:    1234567890,
	// 	UpdatedAt:    1234567890,
	// })
	// queries.CreateUser(ctx, repository.CreateUserParams{
	// 	ID:           ksuid.New().String(),
	// 	Email:        "user@example.com",
	// 	PasswordHash: "hashed_password",
	// 	CreatedAt:    1234567890,
	// 	UpdatedAt:    1234567890,
	// })
	// user, err := queries.GetUser(ctx, "some_user_id")
	// // user, err := queries.GetUser(ctx)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		slog.Info("No user found")
	// 		return
	// 	}
	// 	slog.Error("Failed to get user: " + err.Error())
	// 	return
	// }

	// slog.Info(fmt.Sprintf("User: %s", user.ID))
	// // slog.Info(fmt.Sprintf("User count: %d", len(user)))
	// _ = user
	server := server.NewServer(dbInstance)
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server: " + err.Error())
		return
	}
}
