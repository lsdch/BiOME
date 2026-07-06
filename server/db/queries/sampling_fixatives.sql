-- name: ListUnknownFixativeCodes :many
SELECT c.code::TEXT
FROM unnest(@fixative_codes::text []) AS c(code)
    LEFT JOIN fixatives f ON f.code = c.code
WHERE f.id IS NULL;

-- name: ReplaceFixativesAtSampling :exec
WITH deleted AS (
    DELETE FROM samplings_fixatives
    WHERE sampling_id = @sampling_id::uuid
)
INSERT INTO samplings_fixatives (sampling_id, fixative_id)
SELECT @sampling_id::uuid,
    f.id
FROM fixatives f
WHERE f.code = ANY(@fixative_codes::text []);


-- name: ListFixatives :many
SELECT *
FROM fixatives
ORDER BY name,
    code;

-- name: GetFixativeByCode :one
SELECT *
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
RETURNING *;

-- name: UpdateFixative :one
UPDATE fixatives
SET code = COALESCE(sqlc.narg('code'), code),
    name = COALESCE(sqlc.narg('name'), name),
    description = CASE
        WHEN @set_description::bool THEN sqlc.narg('description')
        ELSE description
    END
WHERE code = @old_code
RETURNING *;

-- name: DeleteFixative :execresult
DELETE FROM fixatives
WHERE code = @code;