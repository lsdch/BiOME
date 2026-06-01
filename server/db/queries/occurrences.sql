-- name: GetOccurrenceByID :one
SELECT o.*,
    sqlc.embed(s),
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
    JOIN samplings s ON s.id = o.sampling_id
    JOIN taxa t ON t.id = o.taxon_id
    JOIN countries c ON c.code = s.site_country_code
WHERE o.id = sqlc.arg(occurrence_id)::uuid;

-- name: GetOccurrenceHabitats :many
SELECT h.*
FROM habitats h
    JOIN samplings_habitats sh ON sh.habitat_id = h.id
    JOIN samplings s ON s.id = sh.sampling_id
    JOIN occurrences o ON o.sampling_id = s.id
WHERE o.id = sqlc.arg(occurrence_id)::uuid;

-- name: GetOccurrenceDatasets :many
SELECT d.*
FROM datasets d
    JOIN occurrences_datasets od ON od.dataset_id = d.id
    JOIN occurrences o ON o.id = od.occurrence_id
WHERE o.id = sqlc.arg(occurrence_id)::uuid; 

-- name: GetSamplingTargetTaxa :many
SELECT t.*
FROM taxa t
    JOIN sampling_target_taxa st ON st.taxon_id = t.id
WHERE st.sampling_id = sqlc.arg(sampling_id)::uuid; 

-- name: ListSamplingSites :many
WITH base_taxa AS (
    SELECT id
    FROM taxa
    WHERE sqlc.arg(taxa_sci_names)::text [] IS NULL
        OR scientific_name = ANY(sqlc.arg(taxa_sci_names))
),
occurring_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE sqlc.arg(taxa_sci_names)::text [] IS NULL
        OR (
            sqlc.arg(include_descendants)::boolean = false
            AND t.id IN (
                SELECT id
                FROM base_taxa
            )
        )
        OR (
            sqlc.arg(include_descendants)::boolean = true
            AND t.id IN (
                SELECT tc.descendant_id
                FROM taxa_closure tc
                    JOIN base_taxa bt ON bt.id = tc.ancestor_id
            )
        )
),
target_base_taxa AS (
    SELECT id
    FROM taxa
    WHERE sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR scientific_name = ANY(sqlc.arg(target_taxa_sci_names))
),
target_sampling_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR (
            sqlc.arg(target_include_descendants)::boolean = false
            AND t.id IN (
                SELECT id
                FROM target_base_taxa
            )
        )
        OR (
            sqlc.arg(target_include_descendants)::boolean = true
            AND t.id IN (
                SELECT tc.descendant_id
                FROM taxa_closure tc
                    JOIN target_base_taxa bt ON bt.id = tc.ancestor_id
            )
        )
),
filtered_sampling_ids AS (
    SELECT DISTINCT st.sampling_id
    FROM sampling_target_taxa st
    WHERE sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR st.taxon_id IN (
            SELECT id
            FROM target_sampling_taxa
        )
),
dataset_filtered_occurrences AS (
    SELECT od.occurrence_id
    FROM occurrences_datasets od
        JOIN datasets d ON d.id = od.dataset_id
    WHERE sqlc.arg(dataset_slugs)::text [] IS NOT NULL
        AND d.slug = ANY(sqlc.arg(dataset_slugs))
),
filtered_occurrences AS (
    SELECT o.*
    FROM occurrences o
        JOIN taxa t ON t.id = o.taxon_id
    WHERE (
            sqlc.arg(taxa_sci_names)::text [] IS NULL
            OR t.id IN (
                SELECT id
                FROM occurring_taxa
            )
        )
        AND (
            sqlc.arg(dataset_slugs)::text [] IS NULL
            OR o.id IN (
                SELECT occurrence_id
                FROM dataset_filtered_occurrences
            )
        )
)
SELECT s.coordinates,
    ARRAY_AGG(DISTINCT s.coordinates_precision) AS coordinates_precision,
    COUNT(DISTINCT s.id)::int AS samplings_count,
    COUNT(DISTINCT o.id)::int AS occurrences_count,
    ARRAY_AGG(DISTINCT s.id) FILTER (
        WHERE s.id IS NOT NULL
    )::uuid [] AS sampling_ids,
    ARRAY_AGG(DISTINCT o.id) FILTER (
        WHERE o.id IS NOT NULL
    )::uuid [] AS occurrence_ids
FROM samplings s
    LEFT JOIN filtered_occurrences o ON o.sampling_id = s.id
    LEFT JOIN samplings_habitats sh ON sh.sampling_id = s.id
    LEFT JOIN habitats h ON h.id = sh.habitat_id
WHERE (
        sqlc.arg(country_code)::text IS NULL
        OR s.site_country_code = sqlc.arg(country_code)
    )
    AND (
        sqlc.arg(habitats)::text [] IS NULL
        OR h.label = ANY(sqlc.arg(habitats))
    )
    AND (
        sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR s.id IN (
            SELECT sampling_id
            FROM filtered_sampling_ids
        )
    )
GROUP BY s.coordinates
HAVING (
        sqlc.arg(min_occurrences)::int IS NULL
        OR COUNT(DISTINCT o.id) > sqlc.arg(min_occurrences)
    );



-- name: OccurrencesGroupsH3 :many
WITH base_taxa AS (
    SELECT id
    FROM taxa
    WHERE sqlc.arg(taxa_sci_names)::text [] IS NULL
        OR scientific_name = ANY(sqlc.arg(taxa_sci_names))
),
occurring_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE sqlc.arg(taxa_sci_names)::text [] IS NULL
        OR (
            sqlc.arg(include_descendants)::boolean = false
            AND t.id IN (
                SELECT id
                FROM base_taxa
            )
        )
        OR (
            sqlc.arg(include_descendants)::boolean = true
            AND t.id IN (
                SELECT tc.descendant_id
                FROM taxa_closure tc
                    JOIN base_taxa bt ON bt.id = tc.ancestor_id
            )
        )
),
target_base_taxa AS (
    SELECT id
    FROM taxa
    WHERE sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR scientific_name = ANY(sqlc.arg(target_taxa_sci_names))
),
target_sampling_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR (
            sqlc.arg(target_include_descendants)::boolean = false
            AND t.id IN (
                SELECT id
                FROM target_base_taxa
            )
        )
        OR (
            sqlc.arg(target_include_descendants)::boolean = true
            AND t.id IN (
                SELECT tc.descendant_id
                FROM taxa_closure tc
                    JOIN target_base_taxa bt ON bt.id = tc.ancestor_id
            )
        )
),
filtered_sampling_ids AS (
    SELECT DISTINCT st.sampling_id
    FROM sampling_target_taxa st
    WHERE sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR st.taxon_id IN (
            SELECT id
            FROM target_sampling_taxa
        )
),
dataset_filtered_occurrences AS (
    SELECT od.occurrence_id
    FROM occurrences_datasets od
        JOIN datasets d ON d.id = od.dataset_id
    WHERE sqlc.arg(dataset_slugs)::text [] IS NOT NULL
        AND d.slug = ANY(sqlc.arg(dataset_slugs))
),
filtered_occurrences AS (
    SELECT o.*
    FROM occurrences o
        JOIN taxa t ON t.id = o.taxon_id
    WHERE (
            sqlc.arg(taxa_sci_names)::text [] IS NULL
            OR t.id IN (
                SELECT id
                FROM occurring_taxa
            )
        )
        AND (
            sqlc.arg(dataset_slugs)::text [] IS NULL
            OR o.id IN (
                SELECT occurrence_id
                FROM dataset_filtered_occurrences
            )
        )
)
SELECT s.h3_res8 AS h3_index,
    COUNT(DISTINCT s.id)::int AS samplings_count,
    COUNT(DISTINCT o.id)::int AS occurrences_count,
    ARRAY_AGG(DISTINCT s.id) FILTER (
        WHERE s.id IS NOT NULL
    )::uuid [] AS sampling_ids,
    ARRAY_AGG(DISTINCT o.id) FILTER (
        WHERE o.id IS NOT NULL
    )::uuid [] AS occurrence_ids
FROM samplings s
    LEFT JOIN filtered_occurrences o ON o.sampling_id = s.id
    LEFT JOIN samplings_habitats sh ON sh.sampling_id = s.id
    LEFT JOIN habitats h ON h.id = sh.habitat_id
WHERE (
        sqlc.arg(country_code)::text IS NULL
        OR s.site_country_code = sqlc.arg(country_code)
    )
    AND (
        sqlc.arg(habitats)::text [] IS NULL
        OR h.label = ANY(sqlc.arg(habitats))
    )
    AND (
        sqlc.arg(target_taxa_sci_names)::text [] IS NULL
        OR s.id IN (
            SELECT sampling_id
            FROM filtered_sampling_ids
        )
    )
GROUP BY s.h3_res8
HAVING (
        sqlc.arg(min_occurrences)::int IS NULL
        OR COUNT(DISTINCT o.id) > sqlc.arg(min_occurrences)
    );