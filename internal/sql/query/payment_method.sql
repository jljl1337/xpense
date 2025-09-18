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

-- name: GetPaymentMethodByID :many
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

-- name: CheckPaymentMethodAccess :one
SELECT
    COUNT(*) > 0 AS can_access
FROM
    payment_method AS pm
LEFT JOIN
    book AS b
ON
    pm.book_id = b.id
WHERE
    pm.id = @payment_method_id AND
    b.user_id = @user_id;