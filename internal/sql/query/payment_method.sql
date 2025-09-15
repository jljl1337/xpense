-- name: CreatePaymentMethod :execrows
INSERT INTO payment_method (
    id,
    book_id,
    name,
    description,
    created_at,
    updated_at
) VALUES (
    @id,
    @book_id,
    @name,
    @description,
    @created_at,
    @updated_at
)
RETURNING
    *;

-- name: GetPaymentMethodsByBookID :many
SELECT
    *
FROM
    payment_method
WHERE
    book_id = @book_id
ORDER BY
    name ASC;

-- name: GetPaymentMethodByID :one
SELECT
    *
FROM
    payment_method
WHERE
    id = @id;

-- name: UpdatePaymentMethodByID :execrows
UPDATE 
    payment_method
SET
    name = @name,
    description = @description,
    updated_at = @updated_at
WHERE
    id = @id;

-- name: DeletePaymentMethodByID :execrows
DELETE FROM
    payment_method
WHERE
    id = @id;