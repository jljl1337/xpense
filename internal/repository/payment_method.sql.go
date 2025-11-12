package repository

import (
	"context"
)

const checkPaymentMethodAccess = `
SELECT
    COUNT(*) > 0 AS can_access
FROM
    payment_method AS pm
LEFT JOIN
    book AS b
ON
    pm.book_id = b.id
WHERE
    pm.id = :id AND
    b.user_id = :user_id
`

type CheckPaymentMethodAccessParams struct {
	PaymentMethodID string `db:"id"`
	UserID          string `db:"user_id"`
}

func (q *Queries) CheckPaymentMethodAccess(ctx context.Context, arg CheckPaymentMethodAccessParams) (bool, error) {
	var canAccess bool
	err := NamedGetContext(ctx, q.db, &canAccess, checkPaymentMethodAccess, arg)
	return canAccess, err
}

const createPaymentMethod = `
INSERT INTO payment_method (
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

type CreatePaymentMethodParams struct {
	ID          string `db:"id"`
	BookID      string `db:"book_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func (q *Queries) CreatePaymentMethod(ctx context.Context, arg CreatePaymentMethodParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, createPaymentMethod, arg)
}

const deletePaymentMethodByID = `
DELETE FROM
    payment_method
WHERE
    id = :id
`

type DeletePaymentMethodByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) DeletePaymentMethodByID(ctx context.Context, id string) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, deletePaymentMethodByID, DeletePaymentMethodByIDParams{ID: id})
}

const getPaymentMethodByID = `
SELECT
    id, book_id, name, description, created_at, updated_at
FROM
    payment_method
WHERE
    id = :id
`

type GetPaymentMethodByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) GetPaymentMethodByID(ctx context.Context, id string) ([]PaymentMethod, error) {
	items := []PaymentMethod{}
	err := NamedSelectContext(ctx, q.db, &items, getPaymentMethodByID, GetPaymentMethodByIDParams{ID: id})
	return items, err
}

const getPaymentMethodsByBookID = `
SELECT
    id, book_id, name, description, created_at, updated_at
FROM
    payment_method
WHERE
    book_id = :book_id
ORDER BY
    name ASC
`

type GetPaymentMethodsByBookIDParams struct {
	BookID string `db:"book_id"`
}

func (q *Queries) GetPaymentMethodsByBookID(ctx context.Context, bookID string) ([]PaymentMethod, error) {
	items := []PaymentMethod{}
	err := NamedSelectContext(ctx, q.db, &items, getPaymentMethodsByBookID, GetPaymentMethodsByBookIDParams{BookID: bookID})
	return items, err
}

const updatePaymentMethodByID = `
UPDATE 
    payment_method
SET
    name = :name,
    description = :description,
    updated_at = :updated_at
WHERE
    id = :id
`

type UpdatePaymentMethodByIDParams struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	UpdatedAt   string `db:"updated_at"`
	ID          string `db:"id"`
}

func (q *Queries) UpdatePaymentMethodByID(ctx context.Context, arg UpdatePaymentMethodByIDParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updatePaymentMethodByID, arg)
}
