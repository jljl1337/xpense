-- name: CreateCategory :execrows
INSERT INTO category (
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