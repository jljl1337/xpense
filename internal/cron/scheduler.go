package cron

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jmoiron/sqlx"

	"github.com/jljl1337/xpense/internal/db"
	"github.com/jljl1337/xpense/internal/env"
	"github.com/jljl1337/xpense/internal/generator"
	"github.com/jljl1337/xpense/internal/repository"
)

func NewScheduler(dbInstance *sqlx.DB) (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	// Database backup job
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

	// Session cleanup job
	if env.SessionCleanupCronSchedule != "" {
		_, err = scheduler.NewJob(
			gocron.CronJob(
				env.SessionCleanupCronSchedule,
				false,
			),
			gocron.NewTask(
				func() {
					slog.Info("Starting session cleanup")

					start := time.Now()

					now := generator.NowISO8601()
					queries := repository.New(dbInstance)
					rows, err := queries.DeleteSessionByExpiresAt(context.Background(), now)
					if err != nil {
						slog.Error("Failed to cleanup expired sessions: " + err.Error())
						return
					}

					slog.Info(fmt.Sprintf("Session cleanup completed in %s, %d sessions deleted", time.Since(start).String(), rows))
				},
			),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create session cleanup cron job: %w", err)
		}
	} else {
		slog.Warn("Session cleanup cron job not scheduled")
	}

	return scheduler, nil
}
