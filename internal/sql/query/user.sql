-- name: CreateUser :execrows
INSERT INTO user (
    id,
    email,
    password_hash,
    created_at,
    updated_at
) VALUES (
    @id,
    @email,
    @password_hash,
    @created_at,
    @updated_at
)
RETURNING
    *;

-- name: GetUserByID :one
SELECT
    *
FROM
    user
WHERE
    id = @id;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    user
WHERE
    email = @email;

-- name: UpdateUser :execrows
UPDATE
    user
SET
    email = @email,
    password_hash = @password_hash,
    updated_at = @updated_at
WHERE
    id = @id;

-- name: DeleteUser :execrows
DELETE FROM
    user
WHERE
    id = @id;