-- name: GetImportBatch :one 
SELECT *
FROM import_batches ib
WHERE ib.id = @import_batch_id;

-- name: ListImportBatches :many
SELECT *
FROM import_batches ib
ORDER BY ib.created_at DESC;

-- name: GetImportBatchForOccurrence :one
SELECT ib.*
FROM import_batches ib
    JOIN occurrences o ON o.import_batch_id = ib.id
WHERE o.id = @occurrence_id;

-- name: DeleteImportBatch :exec
DELETE FROM import_batches
WHERE id = @import_batch_id;

-- name: DeleteOccurrencesFromBatch :exec
WITH deleted_occurrences AS (
    DELETE FROM occurrences o
    WHERE o.import_batch_id = @import_batch_id::ulid
    RETURNING sampling_id
),
affected_events AS (
    SELECT DISTINCT sampling_id
    FROM deleted_occurrences
)
DELETE FROM samplings se
WHERE se.id IN (
        SELECT sampling_id
        FROM affected_events
    )
    AND NOT EXISTS (
        SELECT 1
        FROM occurrences o
        WHERE o.sampling_id = se.id
    );