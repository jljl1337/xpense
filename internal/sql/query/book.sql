-- name: CreateBook :execrows
INSERT INTO book (
    id,
    user_id,
    name,
    description,
    created_at,
    updated_at
) VALUES (
    @id,
    @user_id,
    @name,
    @description,
    @created_at,
    @updated_at
)
RETURNING
    *;

-- name: GetBooksByUserID :many
SELECT
    *
FROM
    book
WHERE
    user_id = @user_id
ORDER BY
    name ASC
LIMIT
    @limit
OFFSET
    @offset;

-- name: GetBookByID :many
SELECT
    *
FROM
    book
WHERE
    id = @id;

-- name: UpdateBookByID :execrows
UPDATE
    book
SET
    name = @name,
    description = @description,
    updated_at = @updated_at
WHERE
    id = @id;

-- name: DeleteBookByID :execrows
DELETE FROM
    book
WHERE
    id = @id;