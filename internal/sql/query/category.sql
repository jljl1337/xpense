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

-- name: GetCategoriesByBookID :many
SELECT
    *
FROM
    category
WHERE
    book_id = @book_id
ORDER BY
    name ASC;

-- name: GetCategoryByID :many
SELECT
    *
FROM
    category
WHERE
    id = @id;

-- name: UpdateCategoryByID :execrows
UPDATE 
    category
SET
    name = @name,
    description = @description,
    updated_at = @updated_at
WHERE
    id = @id;

-- name: DeleteCategoryByID :execrows
DELETE FROM
    category
WHERE
    id = @id;

-- name: CheckCategoryAccess :one
SELECT
    COUNT(*) > 0 AS can_access
FROM
    category AS c
LEFT JOIN
    book AS b
ON
    c.book_id = b.id
WHERE
    c.id = @category_id AND
    b.user_id = @user_id;