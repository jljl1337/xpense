package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jljl1337/xpense/internal/env"
)

func NewDB() (*sql.DB, error) {
	dbPath := env.MustGetString("DB_PATH", "data/live/db/data.db")
	// Create parent directories if they don't exist

	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", "file:"+dbPath+"?_journal=WAL&_foreign_keys=true")
}
