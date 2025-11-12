package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func NamedExecRowsAffectedContext(ctx context.Context, q sqlx.ExtContext, query string, arg any) (int64, error) {
	query, args, err := q.BindNamed(query, arg)
	if err != nil {
		return 0, err
	}
	return ExecRowsAffectedContext(ctx, q, query, args...)
}

func NamedSelectContext(ctx context.Context, q sqlx.ExtContext, dest any, query string, arg any) error {
	query, args, err := q.BindNamed(query, arg)
	if err != nil {
		return err
	}
	return sqlx.SelectContext(ctx, q, dest, query, args...)
}

func NamedGetContext(ctx context.Context, q sqlx.ExtContext, dest any, query string, arg any) error {
	query, args, err := q.BindNamed(query, arg)
	if err != nil {
		return err
	}
	return sqlx.GetContext(ctx, q, dest, query, args...)
}

func ExecRowsAffectedContext(ctx context.Context, q sqlx.ExtContext, query string, args ...any) (int64, error) {
	result, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
