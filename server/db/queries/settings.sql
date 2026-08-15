-- name: GetSettings :one
SELECT *
FROM settings;

-- name: InitSettings :exec
INSERT INTO settings (
        id,
        app_name,
        app_subtitle,
        app_description,
        is_public,
        account_requests_enabled,
        admin_email,
        mail_from_address,
        mail_from_name,
        molecular_data_enabled
    )
VALUES (
        1,
        @app_name,
        @app_subtitle,
        @app_description,
        @is_public,
        @account_requests_enabled,
        @admin_email,
        @mail_from_address,
        @mail_from_name,
        @molecular_data_enabled
    ) ON CONFLICT (id) DO NOTHING;

-- name: UpdateInstanceSettings :one
UPDATE settings
SET app_name = COALESCE(sqlc.narg('app_name'), app_name),
    app_subtitle = CASE
        WHEN @set_app_subtitle::boolean THEN sqlc.narg('app_subtitle')
        ELSE app_subtitle
    END,
    app_description = CASE
        WHEN @set_app_description::boolean THEN sqlc.narg('app_description')
        ELSE app_description
    END,
    is_public = COALESCE(sqlc.narg('is_public'), is_public),
    admin_email = COALESCE(sqlc.narg('admin_email'), admin_email),
    account_requests_enabled = COALESCE(
        sqlc.narg('account_requests_enabled'),
        account_requests_enabled
    ),
    mail_from_address = COALESCE(
        sqlc.narg('mail_from_address'),
        mail_from_address
    ),
    mail_from_name = COALESCE(
        sqlc.narg('mail_from_name'),
        mail_from_name
    ),
    frontpage_message_md = CASE
        WHEN @set_frontpage_message_md::boolean THEN sqlc.narg('frontpage_message_md')
        ELSE frontpage_message_md
    END
WHERE id = 1
RETURNING *;

-- name: SetDashboardMessage :one
UPDATE settings
SET frontpage_message_md = sqlc.narg('frontpage_message_md')
WHERE id = 1
RETURNING frontpage_message_md;