package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

func NewDB(dbPath string) (*sql.DB, error) {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", "file:"+dbPath+"?_journal=WAL&_foreign_keys=true")
}
