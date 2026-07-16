-- name: CheckReadyToMaterialize :one
WITH taxonomy AS (
    SELECT COUNT(*) = 0 AS ready
    FROM taxon_resolution r
    WHERE r.import_id = @import_id
        AND r.resolved_to IS NULL
),
methods AS (
    SELECT COUNT(*) = 0 AS ready
    FROM sampling_methods_resolution r
    WHERE r.import_id = @import_id
        AND r.resolved_method_id IS NULL
)
SELECT taxonomy.ready::bool as taxonomy,
    methods.ready::bool as methods
FROM taxonomy,
    methods;

-- name: MaterializeImportWorkflow :one
INSERT INTO import_batches (
        id,
        workflow_id,
        label,
        description,
        created_by,
        assembled_by
    )
SELECT @batch_ulid,
    w.import_id,
    w.label,
    w.description,
    w.created_by,
    w.assembled_by
FROM import_workflows w
WHERE w.import_id = @import_id::uuid
RETURNING *;

-- name: InitImportWorkflow :one
INSERT INTO import_workflows (label, description, assembled_by, created_by)
VALUES (
        @label,
        @description,
        @assembled_by::TEXT [],
        @created_by::UUID
    )
RETURNING *;

-- name: DeleteImportWorkflow :exec
DELETE FROM import_workflows
WHERE import_id = @import_id::uuid;

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
        id,
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
        @id,
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


-- name: MaterializeSamplings :exec
WITH staging AS (
    SELECT s.*
    FROM samplings_staging s
        JOIN import_workflows w ON w.import_id = s.import_id
        JOIN import_batches b ON b.workflow_id = w.import_id
    WHERE b.id = @import_batch_id::text
),
inserted AS (
    INSERT INTO samplings (
            import_batch_id,
            comments,
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
    SELECT @import_batch_id::text,
        s.sampling_comments,
        s.site_code,
        s.site_name,
        s.site_locality,
        s.site_country_code,
        s.coordinates_precision,
        s.coordinates,
        s.altitude,
        s.event_date,
        s.event_date_precision,
        s.performed_by,
        s.duration,
        s.access_points
    FROM staging s
    RETURNING id
),
staging_with_rn AS (
    SELECT s.import_id,
        s.sampling_hash,
        row_number() OVER (
            ORDER BY s.import_id,
                s.sampling_hash
        ) AS rn
    FROM staging s
),
inserted_with_rn AS (
    SELECT id,
        row_number() OVER () AS rn
    FROM inserted
)
UPDATE import_samplings_occurrences iso
SET materialized_sampling_id = iw.id
FROM inserted_with_rn iw
    JOIN staging_with_rn sw ON iw.rn = sw.rn
WHERE iso.import_id = sw.import_id
    AND iso.sampling_hash = sw.sampling_hash;

-- name: MaterializeOccurrences :exec
INSERT INTO occurrences (
        import_batch_id,
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
SELECT b.id,
    i.id,
    i.occurrence_code,
    i.materialized_sampling_id,
    i.type_status,
    i.occurrence_comments,
    c.taxon_id,
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
    JOIN taxon_resolution r ON r.id = i.taxon_resolution_id
    JOIN taxon_candidates c ON r.resolved_to = c.id
    JOIN import_batches b ON b.workflow_id = i.import_id
WHERE b.id = @batch_id::text;


-- name: CleanUpStagingImport :exec
DELETE FROM import_samplings_occurrences
WHERE import_id = @import_id;


-- name: MaterializeSamplingMethods :exec
INSERT INTO events_sampling_methods (sampling_id, method_id)
SELECT ss.materialized_sampling_id,
    smr.resolved_method_id
FROM samplings_staging ss
    JOIN import_workflows iw ON iw.import_id = ss.import_id
    JOIN import_batches b ON b.workflow_id = iw.import_id -- Si tu as workflow_id dans import_batches :
    -- JOIN import_batches ib ON iw.import_id = ib.workflow_id AND s.import_batch_id = ib.id
    -- Sinon, si tu as import_workflow_id directement dans samplings :
    -- JOIN samplings s ON ss.sampling_hash = s.staging_hash AND ss.import_id = s.import_workflow_id
    JOIN unnest(ss.sampling_methods) AS method_text ON true
    JOIN sampling_methods_resolution smr ON (
        ss.import_id = smr.import_id
        AND smr.input_text = method_text
    )
WHERE b.id = @import_batch_id::text
    AND smr.status = 'selected';