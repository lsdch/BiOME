-- name: GetMethodsResolution :many
SELECT *
FROM sampling_methods_resolution
WHERE import_id = @import_id;

-- name: InitMethodsResolution :many
WITH input_methods AS (
    SELECT DISTINCT UNNEST(sampling_methods) AS input_text
    FROM import_samplings_occurrences
    WHERE import_id = @import_id
),
resolved AS (
    SELECT i.input_text,
        m.id AS resolved_method_id
    FROM input_methods i
        LEFT JOIN sampling_methods m ON lower(unaccent(i.input_text)) = lower(unaccent(m.label))
        OR lower(unaccent(i.input_text)) = lower(unaccent(m.code))
)
INSERT INTO sampling_methods_resolution (
        import_id,
        input_text,
        resolved_method_id,
        status
    )
SELECT @import_id,
    input_text,
    resolved_method_id,
    CASE
        WHEN resolved_method_id IS NOT NULL THEN 'auto'
        ELSE 'pending'
    END::vocab_resolution_status
FROM resolved
RETURNING *;

-- name: ResolveMethod :one 
UPDATE sampling_methods_resolution
SET resolved_method_id = @resolved_method_id,
    status = @status
WHERE import_id = @import_id
    AND input_text = @input_text
RETURNING *;