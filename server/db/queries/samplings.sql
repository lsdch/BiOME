-- name: CreateSampling :one
WITH inserted AS (
    INSERT INTO samplings (
            site_code,
            site_name,
            site_locality,
            site_country_code,
            coordinates,
            coordinates_precision,
            altitude,
            event_date,
            event_date_precision,
            performed_by,
            duration,
            access_points,
            comments
        )
    VALUES (
            @site_code,
            @site_name,
            @site_locality,
            @site_country_code,
            @coordinates,
            @coordinates_precision,
            @altitude,
            @event_date,
            @event_date_precision,
            @performed_by,
            @duration,
            @access_points,
            @comments
        )
)
SELECT sqlc.embed(s),
    sqlc.embed(c)
FROM samplings s
    JOIN countries c ON c.code = s.site_country_code
WHERE s.id = (
        SELECT id
        FROM inserted
    );

-- name: ReplaceSamplingTargetTaxa :exec
WITH deleted AS (
    DELETE FROM sampling_target_taxa
    WHERE sampling_id = @sampling_id::uuid
)
INSERT INTO sampling_target_taxa (sampling_id, taxon_id)
SELECT @sampling_id::uuid,
    t.taxon_id
FROM UNNEST(@taxon_ids::uuid []) AS t(taxon_id);

-- name: GetSampling :one
SELECT sqlc.embed(s),
    sqlc.embed(c)
FROM samplings s
    JOIN countries c ON c.code = s.site_country_code
WHERE s.id = @sampling_id::uuid;

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
    JOIN habitat_groups hg ON hg.id = h.habitat_group_id
WHERE sh.sampling_id = @sampling_id::uuid;

-- name: ListSamplingsH3AtProximity :many
SELECT s.h3_index,
    COUNT(DISTINCT s.id) AS sampling_count,
    COUNT(o.id) AS occurrence_count,
    ST_Distance(
        ST_SetSRID(
            ST_MakePoint(@longitude::real, @latitude::real),
            4326
        )::geography,
        h3_cell_to_geography(s.h3_index::h3index)::geography
    )::integer AS distance_meters
FROM samplings s
    LEFT JOIN occurrences o ON o.sampling_id = s.id
WHERE ST_DWithin(
        s.coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::real, @latitude::real),
            4326
        )::geography,
        @radius_meters::integer
    )
    AND (s.id <> ALL(@exclude_sampling_ids::uuid []))
GROUP BY s.h3_index;

-- name: ListSamplingsAtProximity :many 
SELECT sqlc.embed(s),
    sqlc.embed(c),
    ST_Distance(
        coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::real, @latitude::real),
            4326
        )::geography
    )::integer AS distance_meters
FROM samplings s
    LEFT JOIN countries c ON c.code = s.site_country_code
WHERE ST_DWithin(
        coordinates::geography,
        ST_SetSRID(
            ST_MakePoint(@longitude::real, @latitude::real),
            4326
        )::geography,
        @radius_meters::integer
    )
    AND (
        sqlc.narg('event_date')::date IS NULL
        OR (
            s.event_date_precision IS NOT NULL
            AND s.event_date IS NOT NULL
            AND s.event_date BETWEEN (
                sqlc.narg('event_date')::date - (@date_interval_days::integer) * INTERVAL '1 day'
            )
            AND (
                sqlc.narg('event_date')::date + (@date_interval_days::integer) * INTERVAL '1 day'
            )
        )
    )
    AND (s.id <> ALL(@exclude_sampling_ids::uuid []));

-- name: ListSamplingAccessPoints :many
SELECT ap
FROM (
        SELECT unnest(s.access_points)::TEXT AS ap
        FROM samplings s
    ) sub
GROUP BY ap
ORDER BY ap;