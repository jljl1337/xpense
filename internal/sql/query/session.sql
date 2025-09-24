-- name: CreateSession :execrows
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

-- name: GetSessionByToken :many
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

-- name: UpdateSessionByUserID :execrows
UPDATE
    session
SET
    expires_at = @expires_at,
    updated_at = @updated_at
WHERE
    user_id = @user_id;

-- name: DeleteSessionByExpiresAt :execrows
DELETE FROM
    session
WHERE
    expires_at < @expires_at;