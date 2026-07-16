-- name: DetectBatchSamplingCollisions :many
WITH i AS (
    SELECT *
    FROM import_samplings_occurrences i
    WHERE i.import_id = @import_id
),
s AS (
    -- existing samplings
    SELECT 'existing'::duplicate_source AS source,
        NULL::text AS import_id,
        NULL::int AS row_number,
        s.latitude,
        s.longitude,
        s.event_date,
        s.event_date_precision,
        s.coordinates_precision,
        s.coordinates
    FROM samplings s
    UNION ALL
    -- staging samplings (self comparison)
    SELECT 'staging'::duplicate_source AS source,
        i2.import_id,
        i2.row_number,
        i2.latitude,
        i2.longitude,
        i2.event_date,
        i2.event_date_precision,
        i2.coordinates_precision,
        i2.coordinates
    FROM import_samplings_occurrences i2
    WHERE i2.import_id = @import_id
)
SELECT i.import_id,
    i.sampling_hash,
    i.row_number,
    i.event_date,
    i.event_date_precision,
    i.latitude,
    i.longitude,
    i.coordinates_precision,
    s.source AS duplicate_source,
    s.import_id AS match_import_id,
    s.row_number AS match_row_number,
    s.latitude AS match_latitude,
    s.longitude AS match_longitude,
    s.event_date AS match_event_date,
    s.event_date_precision AS match_event_date_precision,
    s.coordinates_precision AS match_coordinates_precision,
    ST_Distance(
        s.coordinates::geography,
        i.coordinates::geography
    )::int AS distance_meters
FROM i
    JOIN s ON ST_DWithin(
        s.coordinates::geography,
        i.coordinates::geography,
        @radius_meters::int
    )
    AND (
        i.event_date_precision = 'day'
        AND i.event_date IS NOT NULL
        AND s.event_date_precision = 'day'
        AND s.event_date IS NOT NULL
        AND s.event_date BETWEEN (
            i.event_date - make_interval(days => @date_interval_days::int)
        ) AND (
            i.event_date + make_interval(days => @date_interval_days::int)
        )
    );


-- name: DetectBatchOccurrenceCollisions :many
WITH resolved_staging AS (
    -- =========================================================
    -- 1. ENRICH STAGING OCCURRENCES WITH TAXON RESOLUTION
    -- =========================================================
    -- This CTE normalizes imported occurrences by attaching their
    -- resolved taxonomic identity.
    --
    -- Resolution can come from multiple sources:
    --   - internal taxonomy
    --   - GBIF
    --   - manual user resolution
    --
    -- A canonical comparison key is built to unify matching logic
    -- across heterogeneous taxonomic sources.
    --
    SELECT
    i.*,
    c.source,
    c.gbif_id,
    c.taxon_id,
    -- c.staging_id,
    CASE
        WHEN r.source = 'internal' THEN 't:' || r.taxon_id::text
        WHEN r.source = 'gbif' THEN 'g:' || r.gbif_id::text
        WHEN r.source = 'manual' THEN 'm:' || r.staging_id::text
    END AS resolved_taxon_key
    FROM import_samplings_occurrences i
        JOIN taxon_resolution r ON r.import_id = i.import_id
        LEFT JOIN taxon_candidates c ON c.id = r.resolved_to
        AND r.input_name = i.taxon_scientific_name
    WHERE i.import_id = @import_id
),
existing_occurrences AS (
    -- =========================================================
    -- 2. EXISTING OCCURRENCES (PERSISTED DATASET)
    -- =========================================================
    SELECT o.id AS occurrence_id,
        o.taxon_id,
        s.id AS sampling_id,
        s.coordinates,
        s.event_date,
        s.event_date_precision,
        s.coordinates_precision
    FROM occurrences o
        JOIN samplings s ON s.id = o.sampling_id
),
staging_staging_collisions AS (
    -- =========================================================
    -- 3. INTRA-STAGING COLLISIONS (DUPLICATE DETECTION)
    -- =========================================================
    -- Detect duplicates inside the imported dataset.
    --
    -- Two rows collide if:
    --   - same resolved taxon
    --   - spatial proximity (radius)
    --   - temporal proximity (date window)
    --
    SELECT
    a.row_number AS row_a,
    b.row_number AS row_b,
    a.resolved_taxon_key
    FROM resolved_staging a
        JOIN resolved_staging b ON a.resolved_taxon_key = b.resolved_taxon_key
        AND a.row_number < b.row_number
    WHERE ST_DWithin(
            a.coordinates::geography,
            b.coordinates::geography,
            @radius_meters::int
        )
        AND a.event_date IS NOT NULL
        AND b.event_date IS NOT NULL
        AND a.event_date_precision = 'day'
        AND b.event_date_precision = 'day'
        AND a.event_date BETWEEN b.event_date - make_interval(days => @date_interval_days::int) AND b.event_date + make_interval(days => @date_interval_days::int)
),
staging_existing_collisions AS (
    -- =========================================================
    -- 4. STAGING → EXISTING OCCURRENCES COLLISIONS
    -- =========================================================
    -- Detect whether an imported occurrence already exists
    -- in the persisted dataset.
    --
    -- Matching is done on:
    --   - taxonomic identity (resolved_taxon_key → taxon_id)
    --   - spatial proximity
    --   - temporal proximity
    --
    -- NOTE:
    -- Only internal taxonomy is relevant at this stage.
    -- GBIF is excluded by design.
    --
    SELECT
    i.row_number,
    i.resolved_taxon_key,
    o.id AS occurrence_id,
    o.taxon_id AS existing_taxon_id,
    t.name AS existing_taxon_name,
    t.authorship AS existing_taxon_authorship,
    s.id AS sampling_id,
    s.latitude,
    s.longitude,
    s.event_date,
    s.event_date_precision,
    s.coordinates_precision,
    ST_Distance(
        s.coordinates::geography,
        i.coordinates::geography
    )::int AS distance_meters
    FROM resolved_staging i
        JOIN occurrences o ON o.taxon_id = i.taxon_id
        JOIN samplings s ON s.id = o.sampling_id
        JOIN taxa t ON t.id = o.taxon_id
    WHERE ST_DWithin(
            s.coordinates::geography,
            i.coordinates::geography,
            @radius_meters::int
        )
        AND i.event_date_precision = 'day'
        AND i.event_date IS NOT NULL
        AND s.event_date_precision = 'day'
        AND s.event_date IS NOT NULL
        AND s.event_date BETWEEN i.event_date - make_interval(days => @date_interval_days::int) AND i.event_date + make_interval(days => @date_interval_days::int)
) -- =============================================================
-- 5. UNIFIED OUTPUT
-- =============================================================
SELECT 'existing'::duplicate_source AS duplicate_source,
    i.import_id,
    i.row_number,
    i.taxon_name,
    i.taxon_authorship,
    i.latitude,
    i.longitude,
    i.event_date,
    i.event_date_precision,
    i.coordinates_precision,
    i.resolved_taxon_key,
    e.occurrence_id,
    e.existing_taxon_id,
    e.existing_taxon_name as match_taxon_name,
    e.existing_taxon_authorship as match_taxon_authorship,
    e.sampling_id,
    NULL::integer AS match_row_number,
    e.latitude AS match_latitude,
    e.longitude AS match_longitude,
    e.event_date AS match_event_date,
    e.event_date_precision AS match_event_date_precision,
    e.coordinates_precision AS match_coordinates_precision,
    e.distance_meters
FROM staging_existing_collisions e
    JOIN resolved_staging i ON i.row_number = e.row_number
UNION ALL
SELECT 'staging'::duplicate_source AS duplicate_source,
    i.import_id,
    a.row_a AS row_number,
    i.taxon_name,
    i.taxon_authorship,
    i.latitude,
    i.longitude,
    i.event_date,
    i.event_date_precision,
    i.coordinates_precision,
    i.resolved_taxon_key,
    NULL::uuid AS occurrence_id,
    NULL::uuid AS existing_taxon_id,
    b.taxon_name AS match_taxon_name,
    b.taxon_authorship AS match_taxon_authorship,
    NULL::uuid AS sampling_id,
    -- match side (B)
    b.row_number AS match_row_number,
    b.latitude AS match_latitude,
    b.longitude AS match_longitude,
    b.event_date AS match_event_date,
    b.event_date_precision AS match_event_date_precision,
    b.coordinates_precision AS match_coordinates_precision,
    ST_Distance(
        a.coordinates::geography,
        b.coordinates::geography
    )::int AS distance_meters
FROM staging_staging_collisions a
    JOIN resolved_staging i ON i.row_number = a.row_a
    JOIN resolved_staging b ON b.row_number = a.row_b;



-- -- name: DetectBatchOccurrenceCollisions :many
-- SELECT i.import_id,
--     i.row_number,
--     i.event_date,
--     i.event_date_precision,
--     i.latitude,
--     i.longitude,
--     i.coordinates_precision,
--     s.latitude AS sampling_latitude,
--     s.longitude AS sampling_longitude,
--     s.event_date AS sampling_date,
--     s.event_date_precision AS sampling_date_precision,
--     s.coordinates_precision AS sampling_coordinates_precision,
--     ST_Distance(
--         s.coordinates::geography,
--         i.coordinates::geography
--     )::integer AS distance_meters
-- FROM import_samplings_occurrences i
--     JOIN taxon_resolution r ON r.input_name = i.taxon_scientific_name
--     AND r.import_id = i.import_id
--     JOIN samplings s ON ST_DWithin(
--         s.coordinates::geography,
--         i.coordinates::geography,
--         @radius_meters::integer
--     )
--     AND (
--         s.event_date_precision = 'day'
--         AND s.event_date IS NOT NULL
--         AND s.event_date BETWEEN (
--             i.event_date - INTERVAL @date_interval_days::integer || 'day'
--         )
--         AND (
--             i.event_date + INTERVAL @date_interval_days::integer || 'day'
--         )
--     )
--     JOIN occurrences o ON o.sampling_id = s.id
--     AND o.taxon_id = r.taxon_id
-- WHERE i.import_id = $1;