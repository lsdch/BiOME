-- name: GetImportBatch :one 
SELECT *
FROM import_batches ib
WHERE ib.id = @import_batch_id;

-- name: GetImportBatchWithContent :one
SELECT sqlc.embed(ib),
    -- created by user
    sqlc.embed(u),
    -- completed by user
    sqlc.embed(u2),
    COUNT(DISTINCT o.id) AS occurrence_count,
    COUNT(DISTINCT se.id) AS sampling_count
FROM import_batches ib
    LEFT JOIN occurrences o ON o.import_batch_id = ib.id
    LEFT JOIN samplings se ON se.id = o.sampling_id
    JOIN users u ON u.id = ib.created_by
    JOIN users u2 ON u2.id = ib.completed_by
WHERE ib.id = @import_batch_id
    AND ib.status = 'completed'
GROUP BY ib.id,
    u.id,
    u2.id;

-- name: LoadDatasetsForImportBatch :many
SELECT d.*
FROM datasets d
    JOIN datasets_import_batches dib ON dib.dataset_id = d.id
WHERE dib.import_batch_id = @import_batch_id;

-- name: AddImportBatchToDataset :exec
INSERT INTO datasets_import_batches (dataset_id, import_batch_id)
VALUES (
        @dataset_id,
        @import_batch_id
    ) ON CONFLICT DO NOTHING;

-- name: RemoveImportBatchFromDataset :exec
DELETE FROM datasets_import_batches
WHERE dataset_id = @dataset_id
    AND import_batch_id = @import_batch_id;

-- name: ListImportBatches :many
SELECT *
FROM import_batches ib
ORDER BY ib.created_at DESC;

-- name: ListImportBatchesWithContent :many
SELECT sqlc.embed(ib),
    -- created by user
    sqlc.embed(u),
    -- completed by user
    sqlc.embed(u2),
    COUNT(DISTINCT o.id) AS occurrence_count,
    COUNT(DISTINCT se.id) AS sampling_count
FROM import_batches ib
    LEFT JOIN occurrences o ON o.import_batch_id = ib.id
    LEFT JOIN samplings se ON se.id = o.sampling_id
    JOIN users u ON u.id = ib.created_by
    JOIN users u2 ON u2.id = ib.completed_by
WHERE ib.status = 'completed'
GROUP BY ib.id,
    u.id,
    u2.id
ORDER BY ib.created_at DESC;

-- name: ListImportBatchesInDataset :many
SELECT sqlc.embed(ib),
    -- created by user
    sqlc.embed(u),
    -- completed by user
    sqlc.embed(u2),
    COUNT(DISTINCT o.id) AS occurrence_count,
    COUNT(DISTINCT se.id) AS sampling_count
FROM import_batches ib
    LEFT JOIN occurrences o ON o.import_batch_id = ib.id
    LEFT JOIN samplings se ON se.id = o.sampling_id
    JOIN datasets_import_batches dib ON dib.import_batch_id = ib.id
    JOIN users u ON u.id = ib.created_by
    JOIN users u2 ON u2.id = ib.completed_by
WHERE ib.status = 'completed'
    AND dib.dataset_id = @dataset_id
GROUP BY ib.id,
    u.id,
    u2.id
ORDER BY ib.created_at DESC;

-- name: GetImportBatchForOccurrence :one
SELECT ib.*
FROM import_batches ib
    JOIN occurrences o ON o.import_batch_id = ib.id
WHERE o.id = @occurrence_id;

-- name: DeleteImportBatch :one
DELETE FROM import_batches
WHERE id = @import_batch_id
RETURNING *;

-- name: DeleteOccurrencesFromBatch :exec
WITH deleted_occurrences AS (
    DELETE FROM occurrences o
    WHERE o.import_batch_id = @import_batch_id::uuid
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