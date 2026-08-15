-- name: ListDatasets :many
SELECT *
FROM datasets d
ORDER BY d.created_at DESC;

-- name: GetDatasetByID :one
SELECT *
FROM datasets d
WHERE d.id = @dataset_id;

-- name: ListOccurrencesForDataset :many
SELECT sqlc.embed(o),
    sqlc.embed(s),
    sqlc.embed(t)
FROM occurrences o
    JOIN occurrences_datasets od ON od.occurrence_id = o.id
    JOIN samplings_with_country s ON s.id = o.sampling_id
    JOIN countries c ON c.code = s.site_country_code
    JOIN taxa t ON t.id = o.taxon_id
WHERE od.dataset_id = @dataset_id;

-- name: GetDatasetsForOccurrence :many
SELECT d.*
FROM datasets d
    JOIN occurrences_datasets od ON od.dataset_id = d.id
WHERE od.occurrence_id = @occurrence_id;

-- name: AddOccurrenceToDataset :exec
INSERT INTO occurrences_datasets (occurrence_id, dataset_id)
VALUES (
        @occurrence_id::ulid,
        @dataset_id::ulid
    );

-- name: RemoveOccurrenceFromDataset :exec
DELETE FROM occurrences_datasets
WHERE occurrence_id = @occurrence_id::ulid
    AND dataset_id = @dataset_id::ulid;