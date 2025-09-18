-- name: CreateExpense :execrows
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
    @id,
    @book_id,
    @category_id,
    @payment_method_id,
    @date,
    @amount,
    @remark,
    @created_at,
    @updated_at
)
RETURNING
    *;

-- name: GetExpensesByBookID :many
SELECT
    *
FROM
    expense
WHERE
    book_id = @book_id
ORDER BY
    date DESC,
    updated_at DESC
LIMIT
    @limit
OFFSET
    @offset;

-- name: GetExpenseByID :many
SELECT
    *
FROM
    expense
WHERE
    id = @id;

-- name: UpdateExpenseByID :execrows
UPDATE 
    expense
SET
    category_id = @category_id,
    payment_method_id = @payment_method_id,
    date = @date,
    amount = @amount,
    remark = @remark,
    updated_at = @updated_at
WHERE
    id = @id;

-- name: DeleteExpenseByID :execrows
DELETE FROM
    expense
WHERE
    id = @id;

-- name: CheckExpenseAccess :one
SELECT
    COUNT(*) > 0 AS can_access
FROM
    expense AS e
LEFT JOIN
    book AS b
ON
    e.book_id = b.id
WHERE
    e.id = @expense_id AND
    b.user_id = @user_id;

-- name: CountExpensesByCategoryID :one
SELECT
    COUNT(*) AS count
FROM
    expense
WHERE
    category_id = @category_id;

-- name: CountExpensesByPaymentMethodID :one
SELECT
    COUNT(*) AS count
FROM
    expense
WHERE
    payment_method_id = @payment_method_id;