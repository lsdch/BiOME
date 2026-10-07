-- name: GetServiceSettings :one
SELECT *
FROM service_settings
WHERE id = 1;

-- name: InitServiceSettings :exec
INSERT INTO service_settings (id, geoapify_api_key)
VALUES (1, sqlc.narg('geoapify_api_key'))
ON CONFLICT (id) DO NOTHING;

-- name: UpsertServiceSettings :one
INSERT INTO service_settings (id, geoapify_api_key)
VALUES (1, sqlc.narg('geoapify_api_key'))
ON CONFLICT (id) DO UPDATE
SET geoapify_api_key = EXCLUDED.geoapify_api_key
RETURNING *;
