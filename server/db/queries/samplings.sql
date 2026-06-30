-- name: GetSamplingBatch :many
SELECT sqlc.embed(s),
    sqlc.embed(c)
FROM samplings s
    JOIN countries c ON c.code = s.site_country_code
WHERE s.id = ANY(@sampling_ids::uuid []);

-- name: GetOccurrencesAtSamplingsBatch :many
SELECT sqlc.embed(o),
    sqlc.embed(t)
FROM occurrences o
    LEFT JOIN taxa t ON t.id = o.taxon_id
WHERE o.sampling_id = ANY(@sampling_ids::uuid [])
    AND (
        @occurrence_ids::ulid [] IS NULL
        OR o.id = ANY(@occurrence_ids::ulid [])
    );

-- name: GetSamplingMethodsAtEvent :many
SELECT sm.*
FROM sampling_methods sm
    JOIN events_sampling_methods esm ON esm.method_id = sm.id
WHERE esm.sampling_id = @sampling_id::uuid;

-- name: GetSamplingFixativesAtEvent :many
SELECT f.*
FROM fixatives f
    JOIN samplings_fixatives ef ON ef.fixative_id = f.id
WHERE ef.sampling_id = @sampling_id::uuid;

-- name: GetSamplingTargetTaxa :many
SELECT t.*
FROM taxa t
    JOIN sampling_target_taxa st ON st.taxon_id = t.id
WHERE st.sampling_id = @sampling_id::uuid;


-- name: GetHabitatsAtEvent :many
SELECT sqlc.embed(h),
    sqlc.embed(hg)
FROM habitats h
    JOIN samplings_habitats sh ON sh.habitat_id = h.id
    JOIN habitat_groups hg ON hg.id = h.group_id
WHERE sh.sampling_id = @sampling_id::uuid;


-- name: DetectExistingSamplings :many 
SELECT latitude,
    longitude,
    event_date,
    event_date_precision,
    coordinates_precision,
    ST_Distance(
        coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::float8, @latitude::float8),
            4326
        )::geography
    )::integer AS distance_meters
FROM samplings
WHERE ST_DWithin(
        coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::float8, @latitude::float8),
            4326
        )::geography,
        @radius_meters::integer
    )
    AND (
        event_date_precision IS NOT NULL
        AND event_date IS NOT NULL
        AND event_date BETWEEN (
            @event_date::date - INTERVAL @date_interval_days::integer || 'day'
        )
        AND (
            @event_date::date + INTERVAL @date_interval_days::integer || 'day'
        )
    );