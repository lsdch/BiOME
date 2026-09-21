-- name: ListDatasets :many
SELECT sqlc.embed(d),
    sqlc.embed(u),
    COUNT(DISTINCT o.id) AS occurrence_count,
    COUNT(DISTINCT se.id) AS sampling_count,
    COUNT(DISTINCT dib.import_batch_id) AS import_batch_count
FROM datasets d
    JOIN users u ON u.id = d.owner_id
    LEFT JOIN occurrences_datasets od ON od.dataset_id = d.id
    LEFT JOIN datasets_import_batches dib ON dib.dataset_id = d.id
    LEFT JOIN occurrences o ON (
        o.id = od.occurrence_id
        OR o.import_batch_id = dib.import_batch_id
    )
    LEFT JOIN samplings se ON se.id = o.sampling_id
WHERE d.is_public = true
    OR d.owner_id = sqlc.narg('user_id')
GROUP BY d.id,
    u.id
ORDER BY d.created_at DESC;

-- name: GetDatasetByID :one
SELECT sqlc.embed(d),
    sqlc.embed(u)
FROM datasets d
    JOIN users u ON u.id = d.owner_id
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

-- name: CreateDataset :one
INSERT INTO datasets (
        id,
        label,
        slug,
        description,
        pinned,
        owner_id,
        is_public
    )
VALUES (
        @ulid,
        @label,
        @slug,
        @description,
        @pinned,
        @owner_id,
        @is_public
    )
RETURNING *;

-- name: DatasetAddCurator :exec
INSERT INTO datasets_curators (dataset_id, user_id)
VALUES (@dataset_id, @user_id);

-- name: DatasetRemoveCurator :exec
DELETE FROM datasets_curators
WHERE dataset_id = @dataset_id
    AND user_id = @user_id;

-- name: LoadDatasetCurators :many
SELECT sqlc.embed(u),
    dc.dataset_id
FROM users u
    JOIN datasets_curators dc ON dc.user_id = u.id;