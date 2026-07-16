-- name: GetOccurrenceByID :one
SELECT sqlc.embed(o),
    sqlc.embed(s),
    sqlc.embed(t),
    sqlc.embed(c)
FROM occurrences o
    JOIN samplings s ON s.id = o.sampling_id
    JOIN taxa t ON t.id = o.taxon_id
    JOIN countries c ON c.code = s.site_country_code
WHERE o.id = @occurrence_id;

-- name: GetOccurrenceDatasets :many
SELECT d.*
FROM datasets d
    JOIN occurrences_datasets od ON od.dataset_id = d.id
    JOIN occurrences o ON o.id = od.occurrence_id
WHERE o.id = @occurrence_id; 


-- name: GetOccurrenceCodeHistory :many
SELECT h.*
FROM occurrence_code_history h
WHERE h.occurrence_id = @occurrence_id
ORDER BY h.created_at DESC;


-- name: OccurrencesByTaxaOverview :many
SELECT t.id AS id,
    t.name AS name,
    t.authorship AS authorship,
    t.rank AS rank,
    parent.name AS parent_name,
    COUNT(DISTINCT o.id)::int AS occurrences_count,
    COUNT(DISTINCT s.id)::int AS samplings_count,
    ARRAY_AGG(DISTINCT s.id) FILTER (
        WHERE s.id IS NOT NULL
    )::uuid [] AS sampling_ids
FROM taxa t
    LEFT JOIN taxa parent ON parent.id = t.parent_id
    LEFT JOIN occurrences o ON o.taxon_id = t.id
    LEFT JOIN samplings s ON s.id = o.sampling_id
GROUP BY t.id,
    t.name,
    t.authorship,
    t.rank,
    parent.name;


-- name: AddOccurrenceToSampling :one
WITH inserted_occurrence AS (
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
    VALUES (
            @id,
            @code,
            @sampling_id,
            @type_status,
            @comments,
            @taxon_id,
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
        )
)
SELECT sqlc.embed(o),
    sqlc.embed(t)
FROM occurrences o
    JOIN taxa t ON t.id = o.taxon_id
WHERE o.id = @id;


-- name: ListSamplingSites :many
WITH base_taxa AS (
    SELECT id
    FROM taxa
    WHERE @taxa_sci_names::text [] IS NULL
        OR scientific_name = ANY(@taxa_sci_names)
),
occurring_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE @taxa_sci_names::text [] IS NULL
        OR (
            @include_descendants::boolean = false
            AND t.id IN (
                SELECT id
                FROM base_taxa
            )
        )
        OR (
            @include_descendants::boolean = true
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
    WHERE @target_taxa_sci_names::text [] IS NULL
        OR scientific_name = ANY(@target_taxa_sci_names)
),
target_sampling_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE @target_taxa_sci_names::text [] IS NULL
        OR (
            @target_include_descendants::boolean = false
            AND t.id IN (
                SELECT id
                FROM target_base_taxa
            )
        )
        OR (
            @target_include_descendants::boolean = true
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
    WHERE @target_taxa_sci_names::text [] IS NULL
        OR st.taxon_id IN (
            SELECT id
            FROM target_sampling_taxa
        )
),
dataset_filtered_occurrences AS (
    SELECT od.occurrence_id
    FROM occurrences_datasets od
        JOIN datasets d ON d.id = od.dataset_id
    WHERE @dataset_slugs::text [] IS NOT NULL
        AND d.slug = ANY(@dataset_slugs)
),
filtered_occurrences AS (
    SELECT o.*
    FROM occurrences o
        JOIN taxa t ON t.id = o.taxon_id
    WHERE (
            @taxa_sci_names::text [] IS NULL
            OR t.id IN (
                SELECT id
                FROM occurring_taxa
            )
        )
        AND (
            @dataset_slugs::text [] IS NULL
            OR o.id IN (
                SELECT occurrence_id
                FROM dataset_filtered_occurrences
            )
        )
),
filtered_sampling_with_occurrences AS (
    SELECT DISTINCT s.id
    FROM samplings s
        JOIN filtered_occurrences fo ON fo.sampling_id = s.id
),
samplings_without_occurrences AS (
    SELECT s.id
    FROM samplings s
    WHERE NOT EXISTS (
            SELECT 1
            FROM occurrences o
            WHERE o.sampling_id = s.id
        )
),
valid_samplings AS (
    SELECT id
    FROM filtered_sampling_with_occurrences
    UNION ALL
    SELECT id
    FROM samplings_without_occurrences
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
    JOIN valid_samplings vs ON vs.id = s.id
    LEFT JOIN filtered_occurrences o ON o.sampling_id = s.id
    LEFT JOIN samplings_habitats sh ON sh.sampling_id = s.id
    LEFT JOIN habitats h ON h.id = sh.habitat_id
WHERE (
        @country_code::text IS NULL
        OR s.site_country_code = @country_code
    )
    AND (
        @habitats::text [] IS NULL
        OR h.label = ANY(@habitats)
    )
    AND (
        @target_taxa_sci_names::text [] IS NULL
        OR s.id IN (
            SELECT sampling_id
            FROM filtered_sampling_ids
        )
    )
GROUP BY s.coordinates
HAVING (
        @min_occurrences::int IS NULL
        OR COUNT(DISTINCT o.id) > @min_occurrences
    );



-- name: OccurrencesGroupsH3 :many
-- This query aggregates occurrences and samplings by H3 index, 
-- allowing for filtering by taxa, datasets, country, habitats, and minimum occurrences.
WITH base_taxa AS (
    SELECT id
    FROM taxa
    WHERE @taxa_sci_names::text [] IS NULL
        OR scientific_name = ANY(@taxa_sci_names)
),
occurring_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE @taxa_sci_names::text [] IS NULL
        OR (
            @include_descendants::boolean = false
            AND t.id IN (
                SELECT id
                FROM base_taxa
            )
        )
        OR (
            @include_descendants::boolean = true
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
    WHERE @target_taxa_sci_names::text [] IS NULL
        OR scientific_name = ANY(@target_taxa_sci_names)
),
target_sampling_taxa AS (
    SELECT t.id
    FROM taxa t
    WHERE @target_taxa_sci_names::text [] IS NULL
        OR (
            @target_include_descendants::boolean = false
            AND t.id IN (
                SELECT id
                FROM target_base_taxa
            )
        )
        OR (
            @target_include_descendants::boolean = true
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
    WHERE @target_taxa_sci_names::text [] IS NULL
        OR st.taxon_id IN (
            SELECT id
            FROM target_sampling_taxa
        )
),
dataset_filtered_occurrences AS (
    SELECT od.occurrence_id
    FROM occurrences_datasets od
        JOIN datasets d ON d.id = od.dataset_id
    WHERE @dataset_ids::ulid [] IS NOT NULL
        AND d.id = ANY(@dataset_ids)
),
filtered_occurrences AS (
    SELECT o.*
    FROM occurrences o
        JOIN taxa t ON t.id = o.taxon_id
    WHERE (
            @taxa_sci_names::text [] IS NULL
            OR t.id IN (
                SELECT id
                FROM occurring_taxa
            )
        )
        AND (
            @dataset_ids::ulid [] IS NULL
            OR o.id IN (
                SELECT occurrence_id
                FROM dataset_filtered_occurrences
            )
        )
)
SELECT s.h3_index AS h3_index,
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
        @country_codes::text [] IS NULL
        OR s.site_country_code = ANY(@country_codes)
    )
    AND (
        @habitats::text [] IS NULL
        OR h.label = ANY(@habitats)
    )
    AND (
        @target_taxa_sci_names::text [] IS NULL
        OR s.id IN (
            SELECT sampling_id
            FROM filtered_sampling_ids
        )
    )
GROUP BY s.h3_index
HAVING (
        @min_occurrences::int IS NULL
        OR COUNT(DISTINCT o.id) > @min_occurrences
    );


-- name: GetOccurrencesAtH3Index :many
SELECT sqlc.embed(s),
    sqlc.embed(o)
FROM occurrences o
    JOIN samplings s ON o.sampling_id = s.id
WHERE s.h3_index = ANY(
        CASE
            WHEN h3_get_resolution(@h3_index) < 14 THEN h3_to_children(@h3_index, 14)
            WHEN h3_get_resolution(@h3_index) = 14 THEN ARRAY [@h3_index]
            ELSE ARRAY [h3_cell_to_parent(@h3_index, 14)]
        END
    );


-- name: LoadOccurrencesForSamplings :many
SELECT sqlc.embed(o),
    sqlc.embed(t)
FROM occurrences o
    JOIN taxa t ON t.id = o.taxon_id
WHERE o.sampling_id = ANY(@sampling_ids::uuid []);