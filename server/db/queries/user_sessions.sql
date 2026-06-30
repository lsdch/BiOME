-- name: CreateSession :one
INSERT INTO user_sessions (
        id,
        user_id,
        refresh_token_hash,
        expires_at,
        user_agent,
        ip_address
    )
VALUES (
        @session_id,
        @user_id,
        @refresh_token_hash,
        @expires_at,
        @user_agent,
        @ip_address
    )
RETURNING *;

-- name: GetSession :one
SELECT *
FROM user_sessions
WHERE id = @session_id
    AND revoked_at IS NULL
    AND expires_at > now();

-- name: RotateSessionRefreshToken :one
UPDATE user_sessions
SET refresh_token_hash = @new_refresh_token_hash,
    last_used_at = now()
WHERE id = @session_id
    AND refresh_token_hash = @old_refresh_token_hash
    AND revoked_at IS NULL
    AND expires_at > now()
RETURNING *;

-- name: RevokeSession :one
UPDATE user_sessions
SET revoked_at = now()
WHERE id = @session_id
    AND revoked_at IS NULL
RETURNING *;

-- name: RevokeAllSessionsForUser :exec
UPDATE user_sessions
SET revoked_at = now()
WHERE user_id = @user_id
    AND revoked_at IS NULL;