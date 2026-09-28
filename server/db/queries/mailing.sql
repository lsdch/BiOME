-- name: GetMailing :one
SELECT *
FROM mailing
WHERE id = 1;

-- name: UpsertMailing :one
INSERT INTO mailing (
        id,
        mail_from_address,
        mail_from_name,
        smtp_host,
        smtp_port,
        smtp_user,
        smtp_password
    )
VALUES (
        1,
        @mail_from_address,
        @mail_from_name,
        sqlc.narg('smtp_host'),
        sqlc.narg('smtp_port'),
        sqlc.narg('smtp_user'),
        sqlc.narg('smtp_password')
    ) ON CONFLICT (id) DO
UPDATE
SET mail_from_address = EXCLUDED.mail_from_address,
    mail_from_name = EXCLUDED.mail_from_name,
    smtp_host = EXCLUDED.smtp_host,
    smtp_port = EXCLUDED.smtp_port,
    smtp_user = EXCLUDED.smtp_user,
    smtp_password = EXCLUDED.smtp_password,
    last_updated = NOW()
RETURNING *;

-- name: ToggleMailing :exec
UPDATE mailing
SET enabled = @enabled
WHERE id = 1;