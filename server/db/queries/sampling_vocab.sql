-- name: ListSamplingMethods :many
SELECT id,
    code,
    name,
    description
FROM sampling_methods
ORDER BY name,
    code;

-- name: GetSamplingMethodByCode :one
SELECT id,
    code,
    name,
    description
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
RETURNING id,
    code,
    name,
    description;

-- name: UpdateSamplingMethod :one
UPDATE sampling_methods
SET code = COALESCE(sqlc.narg('code'), code),
    name = COALESCE(sqlc.narg('name'), name),
    description = CASE
        WHEN @set_description::bool THEN sqlc.narg('description')
        ELSE description
    END
WHERE code = @old_code
RETURNING id,
    code,
    name,
    description;

-- name: DeleteSamplingMethod :exec
DELETE FROM sampling_methods
WHERE code = @code
RETURNING id,
    code,
    name,
    description;

-- name: ListFixatives :many
SELECT id,
    code,
    name,
    description
FROM fixatives
ORDER BY name,
    code;

-- name: GetFixativeByCode :one
SELECT id,
    code,
    name,
    description
FROM fixatives
WHERE code = @code
LIMIT 1;

-- name: CreateFixative :one
INSERT INTO fixatives (
        code,
        name,
        description
    )
VALUES (
        @code,
        @name,
        @description
    )
RETURNING id,
    code,
    name,
    description;

-- name: UpdateFixative :one
UPDATE fixatives
SET code = COALESCE(sqlc.narg('code'), code),
    name = COALESCE(sqlc.narg('name'), name),
    description = CASE
        WHEN @set_description::bool THEN sqlc.narg('description')
        ELSE description
    END
WHERE code = @old_code
RETURNING id,
    code,
    name,
    description;

-- name: DeleteFixative :execresult
DELETE FROM fixatives
WHERE code = @code;