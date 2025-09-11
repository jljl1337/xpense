-- name: CreateSession :one
INSERT INTO session (
    id,
    user_id,
    token,
    csrf_token,
    expires_at,
    created_at,
    updated_at
) VALUES (
    @id,
    @user_id,
    @token,
    @csrf_token,
    @expires_at,
    @created_at,
    @updated_at
)
RETURNING
    *;

-- name: GetSessionByToken :one
SELECT
    *
FROM
    session
WHERE
    token = @token;

-- name: UpdateSessionByToken :execrows
UPDATE
    session
SET
    expires_at = @expires_at,
    updated_at = @updated_at
WHERE
    token = @token;

-- name: DeleteSession :execrows
DELETE FROM
    session
WHERE
    expires_at < @expires_at;