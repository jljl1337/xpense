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
    @pageSize
OFFSET
    @pageSize * (@page - 1);

-- name: GetExpenseByID :one
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