-- name: CreateInvitation :one
INSERT INTO invitations (
        email,
        invitee_name,
        role,
        message,
        inviter_id,
        expires_at
    )
VALUES (
        @email,
        @invitee_name,
        @role,
        @message,
        @inviter_id,
        @expires_at
    )
RETURNING *;

-- name: CreateInvitationToken :one
INSERT INTO invitation_tokens (invitation_id, token_hash)
VALUES (
        @invitation_id,
        @token_hash
    )
RETURNING *;

-- name: GetInvitationByTokenHash :one
SELECT i.*
FROM invitations i
    JOIN invitation_tokens t ON t.invitation_id = i.id
WHERE t.token_hash = @token_hash
    AND t.consumed = false
    AND i.status = 'pending'
    AND i.expires_at > now()
LIMIT 1;

-- name: CreateUserFromInvitationToken :one
WITH invitation_row AS (
    SELECT i.id,
        i.email,
        i.role
    FROM invitations i
        JOIN invitation_tokens t ON t.invitation_id = i.id
    WHERE t.token_hash = @token_hash
        AND t.consumed = false
        AND i.status = 'pending'
        AND i.expires_at > now() FOR
    UPDATE OF t
),
inserted_user AS (
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
    SELECT @login,
        invitation_row.email,
        @password_hash,
        invitation_row.role,
        @first_name,
        @last_name,
        @organisation,
        @contact,
        @bio,
        now()
    FROM invitation_row
    RETURNING id
),
updated_invitation AS (
    UPDATE invitations i
    SET status = 'redeemed',
        redeemed_at = now()
    FROM inserted_user iu
    WHERE i.id = (
            SELECT id
            FROM invitation_row
        )
    RETURNING i.id
),
updated_token AS (
    UPDATE invitation_tokens t
    SET consumed = true,
        consumed_at = now(),
        consumed_by = iu.id
    FROM inserted_user iu
    WHERE t.token_hash = @token_hash
        AND t.consumed = false
    RETURNING t.id
)
SELECT users.*
FROM users
    JOIN inserted_user ON users.id = inserted_user.id;