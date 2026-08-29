-- name: CheckReadyToMaterialize :one
WITH taxonomy AS (
    SELECT COUNT(*) = 0 AS ready
    FROM taxon_resolution r
    WHERE r.import_id = @import_id
        AND r.resolved_candidate_id IS NULL
),
methods AS (
    SELECT COUNT(*) = 0 AS ready
    FROM sampling_methods_resolution r
    WHERE r.import_id = @import_id
        AND r.resolved_method_id IS NULL
),
fixatives AS (
    SELECT COUNT(*) = 0 AS ready
    FROM sampling_fixatives_resolution r
    WHERE r.import_id = @import_id
        AND r.resolved_fixative_id IS NULL
),
bibliography AS (
    SELECT COUNT(*) = 0 AS ready
    FROM publication_resolution r
    WHERE r.import_id = @import_id
        AND r.status = 'pending'
)
SELECT taxonomy.ready::bool as taxonomy,
    methods.ready::bool as methods,
    fixatives.ready::bool as fixatives,
    bibliography.ready::bool as bibliography
FROM taxonomy,
    methods,
    fixatives,
    bibliography;

-- name: InitImportBatch :one
INSERT INTO import_batches (
        label,
        description,
        assembled_by,
        created_by,
        taxonomic_scope,
        imported_file_name,
        imported_file_size,
        imported_file_hash
    )
VALUES (
        @label,
        @description,
        @assembled_by::TEXT [],
        @created_by::UUID,
        @taxonomic_scope,
        @imported_file_name,
        @imported_file_size,
        @imported_file_hash
    )
RETURNING *;

-- name: SetBatchStatus :exec
UPDATE import_batches
SET status = @status::import_batch_status
WHERE id = @import_id::uuid;

-- name: SetBatchCompleted :exec
UPDATE import_batches
SET status = 'completed',
    completed_at = NOW(),
    completed_by = @completed_by::uuid
WHERE id = @import_id::uuid;

-- name: ListImportBatchs :many
SELECT *
FROM import_batches
ORDER BY created_at DESC;

-- name: GetImportState :one
SELECT *
FROM import_batches
WHERE id = @import_id::uuid;

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

-- name: StageOccurrenceCollections :copyfrom
INSERT INTO collections_staging (occurrence_id, collection_name, vouchers)
VALUES (
        @occurrence_id,
        @collection_name,
        @vouchers
    );

-- name: MaterializeSamplings :exec
WITH batch AS (
    SELECT b.id
    FROM import_batches b
    WHERE b.id = @import_batch_id
),
inserted AS (
    INSERT INTO samplings (
            source_sampling_hash,
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
    SELECT s.sampling_hash,
        b.id,
        s.sampling_comments,
        s.site_code,
        s.site_name,
        s.site_locality,
        COALESCE(s.site_country_code, c.code) AS country_code,
        s.coordinates_precision,
        s.coordinates,
        s.altitude,
        s.event_date,
        s.event_date_precision,
        s.performed_by,
        s.duration,
        s.access_points
    FROM samplings_staging s
        JOIN batch b ON b.id = s.import_id
        LEFT JOIN LATERAL (
            SELECT code
            FROM countries c
            WHERE s.site_country_code IS NULL
                AND s.coordinates IS NOT NULL
                AND ST_Covers(c.geom, s.coordinates)
            LIMIT 1
        ) c ON TRUE
    RETURNING id,
        source_sampling_hash
)
UPDATE import_samplings_occurrences iso
SET materialized_sampling_id = i.id
FROM inserted i
    CROSS JOIN batch b
WHERE iso.import_id = b.id
    AND iso.sampling_hash = i.source_sampling_hash;

-- name: MaterializeOccurrences :exec
WITH inserted AS (
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
        i.generated_code,
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
        JOIN taxon_candidates c ON r.resolved_candidate_id = c.id
        JOIN import_batches b ON b.id = i.import_id
    WHERE b.id = @batch_id
    RETURNING id,
        code
)
INSERT INTO occurrence_code_history (occurrence_id, code)
SELECT inserted.id,
    iso.occurrence_code
FROM inserted
    JOIN import_samplings_occurrences iso ON iso.id = inserted.id
WHERE iso.occurrence_code IS NOT NULL;


-- name: MaterializeCollections :exec
INSERT INTO occurrence_collections (occurrence_id, name, vouchers)
SELECT o.id,
    c.collection_name,
    c.vouchers
FROM collections_staging c
    JOIN occurrences o ON o.id = c.occurrence_id
WHERE o.import_batch_id = @batch_id::uuid;

-- name: CleanUpStagingImport :exec
DELETE FROM import_samplings_occurrences
WHERE import_id = @import_id;


-- name: MaterializeSamplingMethods :exec
INSERT INTO events_sampling_methods (sampling_id, method_id)
SELECT ss.materialized_sampling_id,
    smr.resolved_method_id
FROM samplings_staging ss
    JOIN import_batches b ON b.id = ss.import_id
    JOIN unnest(ss.sampling_methods) AS method_text ON true
    JOIN sampling_methods_resolution smr ON (
        ss.import_id = smr.import_id
        AND smr.input_text = method_text
    )
WHERE b.id = @import_batch_id
    AND smr.status = ANY(
        '{"selected", "auto_resolved"}'::vocab_resolution_status []
    );

-- name: MaterializeSamplingFixatives :exec
INSERT INTO samplings_fixatives (sampling_id, fixative_id)
SELECT ss.materialized_sampling_id,
    sfr.resolved_fixative_id
FROM samplings_staging ss
    JOIN import_batches b ON b.id = ss.import_id
    JOIN unnest(ss.sampling_fixatives) AS fixative_text ON true
    JOIN sampling_fixatives_resolution sfr ON (
        ss.import_id = sfr.import_id
        AND sfr.input_text = fixative_text
    )
WHERE b.id = @import_batch_id
    AND sfr.status = ANY(
        '{"selected", "auto_resolved"}'::vocab_resolution_status []
    );

-- name: MaterializeSamplingTargets :exec
INSERT INTO sampling_target_taxa (sampling_id, taxon_id)
SELECT ss.materialized_sampling_id,
    c.taxon_id
FROM samplings_staging ss
    JOIN sampling_target_resolution r ON (
        ss.import_id = r.import_id
        AND ss.sampling_hash = r.sampling_hash
    )
    JOIN taxon_resolution tr ON tr.id = r.resolution_id
    JOIN taxon_candidates c ON c.id = tr.resolved_candidate_id
WHERE ss.import_id = @import_batch_id;