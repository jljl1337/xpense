package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"

	s "github.com/jljl1337/xpense/internal/sql"
)

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

	appliedMigrations := make(map[string]string)
	for rows.Next() {
		var id string
		var statement string
		if err := rows.Scan(&id, &statement); err != nil {
			return err
		}
		appliedMigrations[id] = statement
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Get the list of migration entryList from the embedded filesystem
	entryList, err := s.MigrationDir.ReadDir("migration")
	if err != nil {
		return err
	}

	// Apply each migration if it hasn't been applied yet
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

		// Skip if the exact same migration has already been applied
		if appliedStatement, ok := appliedMigrations[id]; ok {
			if appliedStatement == statement {
				slog.Debug("Skipping already applied migration: " + id)
				continue
			} else {
				return errors.New("migration conflict with different statements: " + id)
			}
		}

		// Apply the migration within a transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, statement)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx, insertMigration, id, statement)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		// Optionally, update the appliedMigrations map
		appliedMigrations[id] = statement
		slog.Info("Applied migration: " + id)
	}

	return nil
}
