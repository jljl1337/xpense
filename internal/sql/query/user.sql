-- name: CreateUser :execrows
INSERT INTO user (
    id,
    username,
    password_hash,
    created_at,
    updated_at
) VALUES (
    @id,
    @username,
    @password_hash,
    @created_at,
    @updated_at
)
RETURNING
    *;

-- name: GetUserByID :many
SELECT
    *
FROM
    user
WHERE
    id = @id;

-- name: GetUserByUsername :many
SELECT
    *
FROM
    user
WHERE
    username = @username;

-- name: UpdateUserPassword :execrows
UPDATE
    user
SET
    password_hash = @password_hash,
    updated_at = @updated_at
WHERE
    id = @id;

-- name: DeleteUser :execrows
DELETE FROM
    user
WHERE
    id = @id;