-- name: ListUsers :many
SELECT *
FROM users
ORDER BY login ASC;

-- name: CreateUser :one
INSERT INTO users (
        login,
        email,
        password_hash,
        role,
        first_name,
        last_name,
        organisation,
        contact,
        bio,
        email_verified_at
    )
VALUES (
        @login,
        @email,
        @password_hash,
        @role,
        @first_name,
        @last_name,
        @organisation,
        @contact,
        @bio,
        @email_verified_at
    )
RETURNING *;

-- name: SetCurrentUser :exec
SELECT set_config('app.current_user_id', @user_id, true);


-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = @user_id
LIMIT 1;

-- name: GetUserByLoginOrEmail :one
SELECT *
FROM users
WHERE login = @identifier
    OR email = @identifier
LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = sqlc.arg(password_hash)
WHERE id = sqlc.arg(user_id);

-- name: UpdateUserRole :exec
UPDATE users
SET role = sqlc.arg(role)
WHERE id = sqlc.arg(user_id);

-- name: UpdateUserPersonalInfo :one
UPDATE users
SET first_name = CASE
        WHEN sqlc.arg(first_name_set)::boolean THEN sqlc.narg(first_name)
        ELSE first_name
    END,
    last_name = CASE
        WHEN sqlc.arg(last_name_set)::boolean THEN sqlc.narg(last_name)
        ELSE last_name
    END,
    contact = CASE
        WHEN sqlc.arg(contact_set)::boolean THEN sqlc.narg(contact)
        ELSE contact
    END,
    organisation = CASE
        WHEN sqlc.arg(organisation_set)::boolean THEN sqlc.narg(organisation)
        ELSE organisation
    END
WHERE id = sqlc.arg(user_id)
RETURNING *;


-- name: CreateUserAccountRequest :one
INSERT INTO user_account_requests (email, name, motivations, expires_at)
VALUES (
        sqlc.arg(email),
        sqlc.arg(name),
        sqlc.narg(motivations),
        sqlc.arg(expires_at)
    )
RETURNING *;

-- name: CreateUserAccountRequestToken :one
INSERT INTO user_account_request_tokens (user_account_request_id, token_hash)
VALUES (
        sqlc.arg(user_account_request_id),
        sqlc.arg(token_hash)
    )
RETURNING *;

-- name: GetUserAccountRequestByTokenHash :one
SELECT r.*
FROM user_account_requests r
    JOIN user_account_request_tokens t ON t.user_account_request_id = r.id
WHERE t.token_hash = sqlc.arg(token_hash)
    AND t.consumed = false
    AND r.status = 'pending'
    AND r.expires_at > now()
LIMIT 1;

-- name: VerifyUserAccountRequest :one
UPDATE user_account_requests
SET status = 'verified',
    verified_at = now()
WHERE id = sqlc.arg(user_account_request_id)
    AND status = 'pending'
RETURNING *;

-- name: CancelUserAccountRequest :one
UPDATE user_account_requests
SET status = 'cancelled',
    cancelled_at = now()
WHERE id = sqlc.arg(user_account_request_id)
    AND status = 'pending'
RETURNING *;

-- name: ConsumeUserAccountRequestToken :exec
UPDATE user_account_request_tokens
SET consumed = true,
    consumed_at = now()
WHERE token_hash = sqlc.arg(token_hash)
    AND consumed = false;

-- name: CreateUserEmailChangeRequest :one
INSERT INTO user_email_change_requests (user_id, email, expires_at)
VALUES (
        sqlc.arg(user_id),
        sqlc.arg(email),
        sqlc.arg(expires_at)
    )
RETURNING *;

-- name: CreateUserEmailChangeRequestToken :one
INSERT INTO user_email_change_request_tokens (user_email_change_request_id, token_hash)
VALUES (
        sqlc.arg(user_email_change_request_id),
        sqlc.arg(token_hash)
    )
RETURNING *;

-- name: GetUserEmailChangeRequestByTokenHash :one
SELECT r.*
FROM user_email_change_requests r
    JOIN user_email_change_request_tokens t ON t.user_email_change_request_id = r.id
WHERE t.token_hash = sqlc.arg(token_hash)
    AND t.consumed = false
    AND r.status = 'pending'
    AND r.expires_at > now()
LIMIT 1;

-- name: VerifyUserEmailChangeRequest :one
UPDATE user_email_change_requests
SET status = 'verified',
    verified_at = now()
WHERE id = sqlc.arg(user_email_change_request_id)
    AND status = 'pending'
RETURNING *;

-- name: ApplyUserEmailChangeRequest :one
WITH request_row AS (
    SELECT r.id,
        r.user_id,
        r.email
    FROM user_email_change_requests r
    WHERE r.id = sqlc.arg(user_email_change_request_id)
        AND r.status IN ('pending', 'verified')
        AND r.expires_at > now()
    LIMIT 1
), updated_user AS (
    UPDATE users
    SET email = (
            SELECT email
            FROM request_row
        ),
        email_verified_at = now()
    WHERE id = (
            SELECT user_id
            FROM request_row
        )
    RETURNING *
),
updated_request AS (
    UPDATE user_email_change_requests
    SET status = 'applied',
        verified_at = COALESCE(verified_at, now()),
        applied_at = now()
    WHERE id = (
            SELECT id
            FROM request_row
        )
    RETURNING id
)
SELECT *
FROM updated_user;

-- name: ConsumeUserEmailChangeRequestToken :exec
UPDATE user_email_change_request_tokens
SET consumed = true,
    consumed_at = now(),
    consumed_by = sqlc.arg(consumed_by)
WHERE token_hash = sqlc.arg(token_hash)
    AND consumed = false;