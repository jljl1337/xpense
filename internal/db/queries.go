package db

import (
	"database/sql"

	"github.com/jljl1337/xpense/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

func NewRepositoryQueries(db *sql.DB) *repository.Queries {
	return repository.New(db)
}
