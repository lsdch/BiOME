-- name: RefreshOccurrenceCodes :exec
WITH history AS (
    INSERT INTO occurrence_code_history (occurrence_id, code)
    SELECT u.id,
        u.current_code
    FROM occurrence_codes_to_update u
    WHERE u.current_code IS NOT NULL
    RETURNING occurrence_id
)
UPDATE occurrences o
SET code = u.computed_code
FROM occurrence_codes_to_update u
WHERE o.id = u.id
    AND o.code IS DISTINCT
FROM u.computed_code;