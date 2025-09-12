package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jljl1337/xpense/internal/env"
)

func NewDB() (*sql.DB, error) {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(env.DbPath), os.ModePerm); err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", "file:"+env.DbPath+"?_journal=WAL&_foreign_keys=true")
}
