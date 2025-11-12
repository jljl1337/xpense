package repository

import (
	"github.com/jmoiron/sqlx"
)

func New(db sqlx.ExtContext) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db sqlx.ExtContext
}
