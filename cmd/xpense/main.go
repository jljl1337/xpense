package main

import (
	"log/slog"

	"github.com/jljl1337/xpense/internal/db"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/log"
	"github.com/jljl1337/xpense/internal/server"
)

func main() {
	env.SetConstants()

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

	server := server.NewServer(dbInstance)
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server: " + err.Error())
		return
	}
}
