package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"

	s "github.com/jljl1337/xpense/internal/sql"
)

type Migration struct {
	ID        string
	Statement string
}

const createMigrationTable = `
	CREATE TABLE IF NOT EXISTS migration (
		id TEXT PRIMARY KEY,
		statement TEXT NOT NULL
	);
`

const getAppliedMigrations = `
	SELECT
		id,
		statement
	FROM
		migration
	ORDER BY
		id ASC;
`

const insertMigration = `
	INSERT INTO
		migration (
			id,
			statement
		) VALUES (
			?,
			?
		);
`

func Migrate(db *sql.DB) error {
	ctx := context.Background()

	// Create the migrations table if it doesn't exist
	_, err := db.ExecContext(ctx, createMigrationTable)
	if err != nil {
		return err
	}

	// Get the list of applied migrations
	rows, err := db.QueryContext(ctx, getAppliedMigrations)
	if err != nil {
		return err
	}
	defer rows.Close()

	appliedMigrations := make([]Migration, 0)
	for rows.Next() {
		var id string
		var statement string
		if err := rows.Scan(&id, &statement); err != nil {
			return err
		}
		appliedMigrations = append(appliedMigrations, Migration{ID: id, Statement: statement})
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Get all migrations from the embedded filesystem
	entryList, err := s.MigrationDir.ReadDir("migration")
	if err != nil {
		return err
	}

	allMigrations := make([]Migration, 0)

	for _, entry := range entryList {
		// Skip directories
		if entry.IsDir() {
			slog.Warn("Skipping directory in migrations: " + entry.Name())
			continue
		}

		// Get the migration statement
		id := entry.Name()

		statementBytes, err := s.MigrationDir.ReadFile("migration/" + id)
		if err != nil {
			return err
		}
		statement := string(statementBytes)

		allMigrations = append(allMigrations, Migration{ID: id, Statement: statement})
	}

	if len(appliedMigrations) > len(allMigrations) {
		return errors.New("applied migrations are more than the available migrations")
	}

	// Apply the new migrations within a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for i, migration := range allMigrations {
		// Check the corresponding applied migration if it exists
		if i < len(appliedMigrations) {
			appliedMigration := appliedMigrations[i]

			if migration.ID != appliedMigration.ID {
				tx.Rollback()
				return errors.New("migration not found in the applied migrations: " + migration.ID)
			}
			if migration.Statement != appliedMigration.Statement {
				tx.Rollback()
				return errors.New("migration statement does not match the applied migration: " + migration.ID)
			}

			// Migration already applied, skip it
			slog.Debug("Skipping already applied migration: " + migration.ID)
			continue
		}

		// Apply the migration
		_, err = tx.ExecContext(ctx, migration.Statement)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx, insertMigration, migration.ID, migration.Statement)
		if err != nil {
			tx.Rollback()
			return err
		}

		slog.Info("Applied migration: " + migration.ID)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
