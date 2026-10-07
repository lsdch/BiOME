-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
VALUES (@user_id, @token_hash, @expires_at);

-- name: IsPasswordResetTokenValid :one
SELECT EXISTS (
    SELECT 1
    FROM password_reset_tokens t
    JOIN users u ON u.id = t.user_id
    WHERE t.token_hash = @token_hash
        AND t.consumed_at IS NULL
        AND t.expires_at > clock_timestamp()
        AND u.active
);

-- name: GetPasswordResetUserForUpdate :one
SELECT u.*
FROM users u
JOIN password_reset_tokens t ON t.user_id = u.id
WHERE t.token_hash = @token_hash
    AND t.consumed_at IS NULL
    AND t.expires_at > clock_timestamp()
    AND u.active
-- Serialize resets for the same account, including different tokens.
FOR UPDATE OF u;

-- name: ConsumePasswordResetToken :one
UPDATE password_reset_tokens
SET consumed_at = clock_timestamp()
WHERE token_hash = @token_hash
    AND consumed_at IS NULL
    AND expires_at > clock_timestamp()
RETURNING user_id;

-- name: InvalidatePasswordResetTokensForUser :exec
UPDATE password_reset_tokens
SET consumed_at = clock_timestamp()
WHERE user_id = @user_id AND consumed_at IS NULL;
