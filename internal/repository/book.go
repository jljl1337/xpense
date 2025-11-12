package repository

import (
	"context"
)

const createBook = `
INSERT INTO book (
    id,
    user_id,
    name,
    description,
    created_at,
    updated_at
) VALUES (
    :id,
    :user_id,
	:name,
	:description,
	:created_at,
	:updated_at
)
`

type CreateBookParams struct {
	ID          string `db:"id"`
	UserID      string `db:"user_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, createBook, arg)
}

const getBooksCountByUserID = `
SELECT
    COUNT(*) AS count
FROM
    book
WHERE
    user_id = :user_id
`

type GetBooksCountByUserIDParams struct {
	UserID string `db:"user_id"`
}

func (q *Queries) GetBooksCountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := NamedGetContext(ctx, q.db, &count, getBooksCountByUserID, GetBooksCountByUserIDParams{UserID: userID})
	return count, err
}

const getBooksByUserID = `
SELECT
    *
FROM
    book
WHERE
    user_id = :user_id
ORDER BY
    name ASC
LIMIT
    :limit
OFFSET
	:offset
`

type GetBooksByUserIDParams struct {
	UserID string `db:"user_id"`
	Offset int64  `db:"offset"`
	Limit  int64  `db:"limit"`
}

func (q *Queries) GetBooksByUserID(ctx context.Context, arg GetBooksByUserIDParams) ([]Book, error) {
	items := []Book{}
	err := NamedSelectContext(ctx, q.db, &items, getBooksByUserID, arg)
	return items, err
}

const getBookByID = `
SELECT
	*
FROM
    book
WHERE
    id = :id
`

type GetBookByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) GetBookByID(ctx context.Context, id string) ([]Book, error) {
	items := []Book{}
	err := NamedSelectContext(ctx, q.db, &items, getBookByID, GetBookByIDParams{ID: id})
	return items, err
}

const updateBookByID = `
UPDATE
    book
SET
    name = :name,
    description = :description,
    updated_at = :updated_at
WHERE
    id = :id
`

type UpdateBookByIDParams struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	UpdatedAt   string `db:"updated_at"`
	ID          string `db:"id"`
}

func (q *Queries) UpdateBookByID(ctx context.Context, arg UpdateBookByIDParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateBookByID, arg)
}

const deleteBookByID = `
DELETE FROM
    book
WHERE
    id = :id
`

type DeleteBookByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) DeleteBookByID(ctx context.Context, id string) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, deleteBookByID, DeleteBookByIDParams{ID: id})
}

const checkBookAccess = `
SELECT
    COUNT(*) > 0 AS can_access
FROM
    book
WHERE
    id = :book_id AND
    user_id = :user_id
`

type CheckBookAccessParams struct {
	BookID string `db:"book_id"`
	UserID string `db:"user_id"`
}

func (q *Queries) CheckBookAccess(ctx context.Context, arg CheckBookAccessParams) (bool, error) {
	var canAccess bool
	err := NamedGetContext(ctx, q.db, &canAccess, checkBookAccess, arg)
	return canAccess, err
}
