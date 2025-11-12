package repository

import (
	"context"
	"database/sql"
)

const createSession = `
INSERT INTO session (
    id,
    user_id,
    token,
    csrf_token,
    expires_at,
    created_at,
    updated_at
) VALUES (
	:id,
	:user_id,
	:token,
	:csrf_token,
	:expires_at,
	:created_at,
	:updated_at
)
`

type CreateSessionParams struct {
	ID        string         `db:"id"`
	UserID    sql.NullString `db:"user_id"`
	Token     string         `db:"token"`
	CsrfToken string         `db:"csrf_token"`
	ExpiresAt string         `db:"expires_at"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, createSession, arg)
}

const deleteSessionByExpiresAt = `
DELETE FROM
    session
WHERE
    expires_at < :expires_at
`

type DeleteSessionByExpiresAtParams struct {
	ExpiresAt string `db:"expires_at"`
}

func (q *Queries) DeleteSessionByExpiresAt(ctx context.Context, expiresAt string) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, deleteSessionByExpiresAt, DeleteSessionByExpiresAtParams{ExpiresAt: expiresAt})
}

const getSessionByToken = `
SELECT
    id, user_id, token, csrf_token, expires_at, created_at, updated_at
FROM
    session
WHERE
    token = :token
`

type GetSessionByTokenParams struct {
	Token string `db:"token"`
}

func (q *Queries) GetSessionByToken(ctx context.Context, token string) ([]Session, error) {
	items := []Session{}
	err := NamedSelectContext(ctx, q.db, &items, getSessionByToken, GetSessionByTokenParams{Token: token})
	return items, err
}

const updateSessionByToken = `
UPDATE
    session
SET
    expires_at = :expires_at,
    updated_at = :updated_at
WHERE
    token = :token
`

type UpdateSessionByTokenParams struct {
	ExpiresAt string `db:"expires_at"`
	UpdatedAt string `db:"updated_at"`
	Token     string `db:"token"`
}

func (q *Queries) UpdateSessionByToken(ctx context.Context, arg UpdateSessionByTokenParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateSessionByToken, arg)
}

const updateSessionByUserID = `
UPDATE
    session
SET
    expires_at = :expires_at,
    updated_at = :updated_at
WHERE
    user_id = :user_id AND
    expires_at > :expires_at
`

type UpdateSessionByUserIDParams struct {
	ExpiresAt string         `db:"expires_at"`
	UpdatedAt string         `db:"updated_at"`
	UserID    sql.NullString `db:"user_id"`
}

func (q *Queries) UpdateSessionByUserID(ctx context.Context, arg UpdateSessionByUserIDParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateSessionByUserID, arg)
}
