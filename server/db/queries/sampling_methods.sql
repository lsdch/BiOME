-- name: ListUnknownSamplingMethodCodes :many
SELECT c.code::TEXT
FROM unnest(@method_codes::text []) AS c(code)
    LEFT JOIN sampling_methods sm ON sm.code = c.code
WHERE sm.id IS NULL;

-- name: ReplaceMethodsAtSampling :exec
WITH deleted AS (
    DELETE FROM events_sampling_methods
    WHERE sampling_id = @sampling_id::uuid
)
INSERT INTO events_sampling_methods (sampling_id, method_id)
SELECT @sampling_id::uuid,
    sm.id
FROM sampling_methods sm
WHERE sm.code = ANY(@method_codes::text []);

-- name: ListSamplingMethods :many
SELECT *
FROM sampling_methods
ORDER BY name,
    code;

-- name: GetSamplingMethodByCode :one
SELECT *
FROM sampling_methods
WHERE code = @code
LIMIT 1;

-- name: CreateSamplingMethod :one
INSERT INTO sampling_methods (
        code,
        name,
        description
    )
VALUES (
        @code,
        @name,
        @description
    )
RETURNING *;

-- name: UpdateSamplingMethod :one
UPDATE sampling_methods
SET code = COALESCE(sqlc.narg('code'), code),
    name = COALESCE(sqlc.narg('name'), name),
    description = CASE
        WHEN @set_description::bool THEN sqlc.narg('description')
        ELSE description
    END
WHERE code = @old_code
RETURNING *;

-- name: DeleteSamplingMethod :exec
DELETE FROM sampling_methods
WHERE code = @code
RETURNING *;