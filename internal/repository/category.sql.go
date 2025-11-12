package repository

import (
	"context"
)

const checkCategoryAccess = `
SELECT
    COUNT(*) > 0 AS can_access
FROM
    category AS c
LEFT JOIN
    book AS b
ON
    c.book_id = b.id
WHERE
    c.id = :id AND
    b.user_id = :user_id
`

type CheckCategoryAccessParams struct {
	CategoryID string `db:"id"`
	UserID     string `db:"user_id"`
}

func (q *Queries) CheckCategoryAccess(ctx context.Context, arg CheckCategoryAccessParams) (bool, error) {
	var canAccess bool
	err := NamedGetContext(ctx, q.db, &canAccess, checkCategoryAccess, arg)
	return canAccess, err
}

const createCategory = `
INSERT INTO category (
    id,
    book_id,
    name,
    description,
    created_at,
    updated_at
) VALUES (
    :id,
    :book_id,
	:name,
	:description,
	:created_at,
	:updated_at
)
`

type CreateCategoryParams struct {
	ID          string `db:"id"`
	BookID      string `db:"book_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, createCategory, arg)
}

const deleteCategoryByID = `
DELETE FROM
    category
WHERE
    id = :id
`

type DeleteCategoryByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) DeleteCategoryByID(ctx context.Context, id string) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, deleteCategoryByID, DeleteCategoryByIDParams{ID: id})
}

const getCategoriesByBookID = `
SELECT
    id, book_id, name, description, created_at, updated_at
FROM
    category
WHERE
    book_id = :book_id
ORDER BY
    name ASC
`

type GetCategoriesByBookIDParams struct {
	BookID string `db:"book_id"`
}

func (q *Queries) GetCategoriesByBookID(ctx context.Context, bookID string) ([]Category, error) {
	items := []Category{}
	err := NamedSelectContext(ctx, q.db, &items, getCategoriesByBookID, GetCategoriesByBookIDParams{BookID: bookID})
	return items, err
}

const getCategoryByID = `
SELECT
    id, book_id, name, description, created_at, updated_at
FROM
    category
WHERE
    id = :id
`

type GetCategoryByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) GetCategoryByID(ctx context.Context, id string) ([]Category, error) {
	items := []Category{}
	err := NamedSelectContext(ctx, q.db, &items, getCategoryByID, GetCategoryByIDParams{ID: id})
	return items, err
}

const updateCategoryByID = `
UPDATE 
    category
SET
    name = :name,
    description = :description,
    updated_at = :updated_at
WHERE
    id = :id
`

type UpdateCategoryByIDParams struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	UpdatedAt   string `db:"updated_at"`
	ID          string `db:"id"`
}

func (q *Queries) UpdateCategoryByID(ctx context.Context, arg UpdateCategoryByIDParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateCategoryByID, arg)
}
