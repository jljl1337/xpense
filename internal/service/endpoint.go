package service

import "github.com/jmoiron/sqlx"

type EndpointService struct {
	db *sqlx.DB
}

func NewEndpointService(db *sqlx.DB) *EndpointService {
	return &EndpointService{
		db: db,
	}
}
