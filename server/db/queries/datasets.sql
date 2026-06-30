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
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
    JOIN occurrences_datasets od ON od.occurrence_id = o.id
    JOIN samplings s ON s.id = o.sampling_id
    JOIN countries c ON c.code = s.site_country_code
    JOIN taxa t ON t.id = o.taxon_id
WHERE od.dataset_id = @dataset_id;

-- name: GetDatasetsForOccurrence :many
SELECT d.*
FROM datasets d
    JOIN occurrences_datasets od ON od.dataset_id = d.id
WHERE od.occurrence_id = @occurrence_id;