-- name: InitBatchImport :one
INSERT INTO import_workflows (label)
VALUES (@label::TEXT)
RETURNING *;

-- name: ListImportWorkflows :many
SELECT *
FROM import_workflows
ORDER BY created_at DESC;

-- name: GetImportState :one
SELECT *
FROM import_workflows
WHERE import_id = @import_id::uuid;

-- name: CopyImportStaging :copyfrom
INSERT INTO import_samplings_occurrences (
        import_id,
        sampling_hash,
        row_number,
        sampling_comments,
        site_code,
        site_name,
        site_locality,
        site_country_code,
        coordinates_precision,
        longitude,
        latitude,
        altitude,
        event_date,
        event_date_precision,
        performed_by,
        duration,
        access_points,
        sampling_targets,
        sampling_fixatives,
        sampling_methods,
        habitats,
        occurrence_code,
        type_status,
        occurrence_comments,
        taxon_name,
        taxon_rank,
        taxon_authorship,
        verbatim_identification,
        identified_by,
        identification_date,
        identification_date_precision,
        identification_confer,
        identification_addendum,
        content_description,
        quantity_exact,
        quantity_lower,
        quantity_upper,
        sources
    )
VALUES (
        @import_id,
        @sampling_hash,
        @row_number,
        @sampling_comments,
        @site_code,
        @site_name,
        @site_locality,
        @site_country_code,
        @coordinates_precision,
        @longitude,
        @latitude,
        @altitude,
        @event_date,
        @event_date_precision,
        @performed_by,
        @duration,
        @access_points,
        @sampling_targets,
        @sampling_fixatives,
        @sampling_methods,
        @habitats,
        @occurrence_code,
        @type_status,
        @occurrence_comments,
        @taxon_name,
        @taxon_rank,
        @taxon_authorship,
        @verbatim_identification,
        @identified_by,
        @identification_date,
        @identification_date_precision,
        @identification_confer,
        @identification_addendum,
        @content_description,
        @quantity_exact,
        @quantity_lower,
        @quantity_upper,
        @sources
    );

-- name: UpsertSamplingsFromStaging :many
INSERT INTO samplings (
        sampling_hash,
        notes,
        site_code,
        site_name,
        site_locality,
        site_country_code,
        coordinates_precision,
        coordinates,
        altitude,
        event_date,
        event_date_precision,
        performed_by,
        duration,
        access_points
    )
SELECT DISTINCT sampling_hash,
    notes,
    site_code,
    site_name,
    site_locality,
    site_country_code,
    coordinates_precision,
    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326),
    altitude,
    event_date,
    event_date_precision,
    performed_by,
    duration,
    access_points
FROM import_samplings_occurrences
WHERE import_id = $1 ON CONFLICT (sampling_hash) DO
UPDATE
SET notes = EXCLUDED.notes
RETURNING id;


-- name: InsertOccurrencesFromStaging :exec
INSERT INTO occurrences (
        id,
        code,
        sampling_id,
        type_status,
        comments,
        taxon_id,
        verbatim_identification,
        identified_by,
        identification_date,
        identification_date_precision,
        identification_confer,
        identification_addendum,
        content_description,
        quantity_exact,
        quantity_lower,
        quantity_upper,
        sources
    )
SELECT i.occurrence_id,
    i.occurrence_code,
    s.id,
    i.type_status,
    i.comments,
    r.taxon_id,
    i.verbatim_identification,
    i.identified_by,
    i.identification_date,
    i.identification_date_precision,
    i.identification_confer,
    i.identification_addendum,
    i.content_description,
    i.quantity_exact,
    i.quantity_lower,
    i.quantity_upper,
    i.sources
FROM import_samplings_occurrences i
    JOIN samplings s ON s.sampling_hash = i.sampling_hash
    JOIN taxon_candidates r ON r.import_id = i.import_id
    AND r.input_name = i.taxon_scientific_name
WHERE i.import_id = $1;

-- name: CleanUpStagingImport :exec
DELETE FROM import_samplings_occurrences
WHERE import_id = $1;