package repository

import (
	"context"
)

const createUser = `
INSERT INTO user (
    id,
    username,
    password_hash,
    created_at,
    updated_at
) VALUES (
	:id,
	:username,
	:password_hash,
	:created_at,
	:updated_at
)
`

type CreateUserParams struct {
	ID           string `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, createUser, arg)
}

const deleteUser = `
DELETE FROM
    user
WHERE
    id = :id
`

type DeleteUserParams struct {
	ID string `db:"id"`
}

func (q *Queries) DeleteUser(ctx context.Context, id string) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, deleteUser, DeleteUserParams{ID: id})
}

const getUserByID = `
SELECT
    id, username, password_hash, created_at, updated_at
FROM
    user
WHERE
    id = :id
`

type GetUserByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) GetUserByID(ctx context.Context, id string) ([]User, error) {
	items := []User{}
	err := NamedSelectContext(ctx, q.db, &items, getUserByID, GetUserByIDParams{ID: id})
	return items, err
}

const getUserByUsername = `
SELECT
    id, username, password_hash, created_at, updated_at
FROM
    user
WHERE
    username = :username
`

type GetUserByUsernameParams struct {
	Username string `db:"username"`
}

func (q *Queries) GetUserByUsername(ctx context.Context, username string) ([]User, error) {
	items := []User{}
	err := NamedSelectContext(ctx, q.db, &items, getUserByUsername, GetUserByUsernameParams{Username: username})
	return items, err
}

const updateUserPassword = `
UPDATE
    user
SET
    password_hash = :password_hash,
    updated_at = :updated_at
WHERE
    id = :id
`

type UpdateUserPasswordParams struct {
	PasswordHash string `db:"password_hash"`
	UpdatedAt    string `db:"updated_at"`
	ID           string `db:"id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateUserPassword, arg)
}

const updateUserUsername = `
UPDATE
    user
SET
    username = :username,
    updated_at = :updated_at
WHERE
    id = :id
`

type UpdateUserUsernameParams struct {
	Username  string `db:"username"`
	UpdatedAt string `db:"updated_at"`
	ID        string `db:"id"`
}

func (q *Queries) UpdateUserUsername(ctx context.Context, arg UpdateUserUsernameParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateUserUsername, arg)
}
