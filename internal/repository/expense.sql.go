package repository

import (
	"context"
)

const checkExpenseAccess = `
SELECT
    COUNT(*) > 0 AS can_access
FROM
    expense AS e
LEFT JOIN
    book AS b
ON
    e.book_id = b.id
WHERE
    e.id = :id AND
    b.user_id = :user_id
`

type CheckExpenseAccessParams struct {
	ExpenseID string `db:"id"`
	UserID    string `db:"user_id"`
}

func (q *Queries) CheckExpenseAccess(ctx context.Context, arg CheckExpenseAccessParams) (bool, error) {
	var canAccess bool
	err := NamedGetContext(ctx, q.db, &canAccess, checkExpenseAccess, arg)
	return canAccess, err
}

const countExpensesByCategoryID = `
SELECT
    COUNT(*) AS count
FROM
    expense
WHERE
    category_id = :category_id
`

type CountExpensesByCategoryIDParams struct {
	CategoryID string `db:"category_id"`
}

func (q *Queries) CountExpensesByCategoryID(ctx context.Context, categoryID string) (int64, error) {
	var count int64
	err := NamedGetContext(ctx, q.db, &count, countExpensesByCategoryID, CountExpensesByCategoryIDParams{CategoryID: categoryID})
	return count, err
}

const countExpensesByPaymentMethodID = `
SELECT
    COUNT(*) AS count
FROM
    expense
WHERE
    payment_method_id = :payment_method_id
`

type CountExpensesByPaymentMethodIDParams struct {
	PaymentMethodID string `db:"payment_method_id"`
}

func (q *Queries) CountExpensesByPaymentMethodID(ctx context.Context, paymentMethodID string) (int64, error) {
	var count int64
	err := NamedGetContext(ctx, q.db, &count, countExpensesByPaymentMethodID, CountExpensesByPaymentMethodIDParams{PaymentMethodID: paymentMethodID})
	return count, err
}

const createExpense = `
INSERT INTO expense (
    id,
    book_id,
    category_id,
    payment_method_id,
    date,
    amount,
    remark,
    created_at,
    updated_at
) VALUES (
    :id,
	:book_id,
	:category_id,
	:payment_method_id,
	:date,
	:amount,
	:remark,
	:created_at,
	:updated_at
)
`

type CreateExpenseParams struct {
	ID              string  `db:"id"`
	BookID          string  `db:"book_id"`
	CategoryID      string  `db:"category_id"`
	PaymentMethodID string  `db:"payment_method_id"`
	Date            string  `db:"date"`
	Amount          float64 `db:"amount"`
	Remark          string  `db:"remark"`
	CreatedAt       string  `db:"created_at"`
	UpdatedAt       string  `db:"updated_at"`
}

func (q *Queries) CreateExpense(ctx context.Context, arg CreateExpenseParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, createExpense, arg)
}

const deleteExpenseByID = `
DELETE FROM
    expense
WHERE
    id = :id
`

type DeleteExpenseByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) DeleteExpenseByID(ctx context.Context, id string) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, deleteExpenseByID, DeleteExpenseByIDParams{ID: id})
}

const getExpenseByID = `
SELECT
	*
FROM
    expense
WHERE
    id = :id
`

type GetExpenseByIDParams struct {
	ID string `db:"id"`
}

func (q *Queries) GetExpenseByID(ctx context.Context, id string) ([]Expense, error) {
	items := []Expense{}
	err := NamedSelectContext(ctx, q.db, &items, getExpenseByID, GetExpenseByIDParams{ID: id})
	return items, err
}

const getExpenseCountByBookID = `
SELECT
    COUNT(*) AS count
FROM
    expense
WHERE
    book_id = :book_id AND
    (category_id = :category_id OR :category_id = '') AND
    (payment_method_id = :payment_method_id OR :payment_method_id = '') AND
    (INSTR(remark, :remark) > 0 OR :remark = '')
`

type GetExpenseCountByBookIDParams struct {
	BookID          string `db:"book_id"`
	CategoryID      string `db:"category_id"`
	PaymentMethodID string `db:"payment_method_id"`
	Remark          string `db:"remark"`
}

func (q *Queries) GetExpenseCountByBookID(ctx context.Context, arg GetExpenseCountByBookIDParams) (int64, error) {
	var count int64
	err := NamedGetContext(ctx, q.db, &count, getExpenseCountByBookID, arg)
	return count, err
}

const getExpensesByBookID = `
SELECT
    *
FROM
    expense
WHERE
    book_id = :book_id AND
    (category_id = :category_id OR :category_id = '') AND
    (payment_method_id = :payment_method_id OR :payment_method_id = '') AND
    (INSTR(remark, :remark) > 0 OR :remark = '')
ORDER BY
    date DESC,
    updated_at DESC
LIMIT
    :limit
OFFSET
    :offset
`

type GetExpensesByBookIDParams struct {
	BookID          string `db:"book_id"`
	CategoryID      string `db:"category_id"`
	PaymentMethodID string `db:"payment_method_id"`
	Remark          string `db:"remark"`
	Offset          int64  `db:"offset"`
	Limit           int64  `db:"limit"`
}

func (q *Queries) GetExpensesByBookID(ctx context.Context, arg GetExpensesByBookIDParams) ([]Expense, error) {
	items := []Expense{}
	err := NamedSelectContext(ctx, q.db, &items, getExpensesByBookID, arg)
	return items, err
}

const updateExpenseByID = `
UPDATE 
    expense
SET
    category_id = :category_id,
    payment_method_id = :payment_method_id,
    date = :date,
    amount = :amount,
    remark = :remark,
    updated_at = :updated_at
WHERE
    id = :id
`

type UpdateExpenseByIDParams struct {
	CategoryID      string  `db:"category_id"`
	PaymentMethodID string  `db:"payment_method_id"`
	Date            string  `db:"date"`
	Amount          float64 `db:"amount"`
	Remark          string  `db:"remark"`
	UpdatedAt       string  `db:"updated_at"`
	ID              string  `db:"id"`
}

func (q *Queries) UpdateExpenseByID(ctx context.Context, arg UpdateExpenseByIDParams) (int64, error) {
	return NamedExecRowsAffectedContext(ctx, q.db, updateExpenseByID, arg)
}
