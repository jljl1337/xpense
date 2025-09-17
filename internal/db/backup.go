package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

// BackupToFile backs up a database to a file
func BackupToFile(srcDB *sql.DB, backupPath string) error {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(backupPath), os.ModePerm); err != nil {
		return err
	}

	// Open destination database
	destDB, err := sql.Open("sqlite3", backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup database: %w", err)
	}
	defer destDB.Close()

	// Perform backup
	return backup(destDB, srcDB)
}

// backup performs a complete backup from srcDb to destDb
func backup(destDb, srcDb *sql.DB) error {
	// Get raw connections
	destConn, err := destDb.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get destination connection: %w", err)
	}
	defer destConn.Close()

	srcConn, err := srcDb.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get source connection: %w", err)
	}
	defer srcConn.Close()

	// Perform backup using raw connections
	return destConn.Raw(func(destConn any) error {
		return srcConn.Raw(func(srcConn any) error {
			// Convert to SQLite connections
			destSQLiteConn, ok := destConn.(*sqlite3.SQLiteConn)
			if !ok {
				return fmt.Errorf("can't convert destination connection to SQLiteConn")
			}

			srcSQLiteConn, ok := srcConn.(*sqlite3.SQLiteConn)
			if !ok {
				return fmt.Errorf("can't convert source connection to SQLiteConn")
			}

			// Initialize backup
			b, err := destSQLiteConn.Backup("main", srcSQLiteConn, "main")
			if err != nil {
				return fmt.Errorf("error initializing SQLite backup: %w", err)
			}

			// Perform backup in one step (-1 means copy entire database)
			done, err := b.Step(-1)
			if err != nil {
				return fmt.Errorf("error in stepping backup: %w", err)
			}
			if !done {
				return fmt.Errorf("backup not completed in one step")
			}

			// Finish backup
			if err := b.Finish(); err != nil {
				return fmt.Errorf("error finishing backup: %w", err)
			}

			return nil
		})
	})
}
