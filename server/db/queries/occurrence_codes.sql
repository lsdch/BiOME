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

-- name: GenerateCodesStaging :exec
UPDATE import_samplings_occurrences
SET generated_code = o.computed_code
FROM (
        SELECT o.id,
            generate_occurrence_code(
                t.name::text,
                o.site_code,
                o.latitude,
                o.longitude,
                o.event_date,
                o.event_date_precision
            ) AS computed_code
        FROM import_samplings_occurrences o
            JOIN taxon_resolution r ON r.id = o.taxon_resolution_id
            JOIN taxon_candidates c ON r.resolved_candidate_id = c.id
            JOIN taxa t ON t.id = c.taxon_id
            JOIN import_batches b ON b.id = o.import_id
        WHERE b.id = @batch_id
    ) AS o(id, computed_code)
WHERE import_samplings_occurrences.id = o.id;

-- name: CheckStagingCodesGenerated :many
SELECT *
FROM import_samplings_occurrences o
    JOIN import_batches b ON b.id = o.import_id
WHERE o.generated_code IS NULL
    AND b.id = @batch_id;