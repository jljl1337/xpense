package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

func NewDB(dbPath, dbBusyTimeout string) (*sql.DB, error) {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
		return nil, err
	}

	dsn := "file:" + dbPath
	dsn = dsn + "?_journal=WAL"
	dsn = dsn + "&_foreign_keys=true"
	dsn = dsn + "&_busy_timeout=" + dbBusyTimeout
	return sql.Open("sqlite3", dsn)
}
