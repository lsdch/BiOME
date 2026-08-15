-- name: InsertGBIFBatch :batchexec
INSERT INTO gbif_staging (
        key,
        parent,
        parent_key,
        canonical_name,
        scientific_name,
        status,
        rank,
        name_type,
        kingdom_key,
        phylum_key,
        class_key,
        order_key,
        family_key,
        genus_key,
        species_key,
        higher_taxon_keys,
        higher_taxon_names,
        authorship,
        num_descendants,
        accepted_key,
        accepted_name
    )
VALUES (
        @key,
        @parent,
        @parent_key,
        @canonical_name,
        @scientific_name,
        @status,
        @rank,
        @name_type,
        @kingdom_key,
        @phylum_key,
        @class_key,
        @order_key,
        @family_key,
        @genus_key,
        @species_key,
        @higher_taxon_keys,
        @higher_taxon_names,
        @authorship,
        @num_descendants,
        @accepted_key,
        @accepted_name
    ) ON CONFLICT (key) DO NOTHING;

-- name: CleanUpGBIFDependencies :exec
DELETE FROM gbif_dependencies d
WHERE d.import_id = @import_id;

-- name: ExpandGBIFDependencies :exec
-- This query expands the list of GBIF keys that need to be resolved for a given import batch. 
-- It includes direct GBIF IDs, accepted keys, and higher taxonomy keys.
WITH candidates AS (
    SELECT DISTINCT c.import_id,
        c.gbif_id
    FROM taxon_candidates c
        JOIN gbif_staging g ON g.key = c.gbif_id
    WHERE c.import_id = @import_id
        AND c.gbif_id IS NOT NULL
),
expanded AS (
    -- 1. direct gbif_id
    SELECT r.import_id,
        r.gbif_id AS key,
        r.gbif_id AS from_key
    FROM candidates r
    UNION
    -- 2. accepted key
    SELECT r.import_id,
        g.accepted_key AS key,
        r.gbif_id AS from_key
    FROM candidates r
        JOIN gbif_staging g ON g.key = r.gbif_id
    WHERE g.accepted_key IS NOT NULL
    UNION
    -- 3. higher taxonomy expansion
    SELECT r.import_id,
        h.key,
        r.gbif_id AS from_key
    FROM candidates r
        JOIN gbif_staging g ON g.key = r.gbif_id
        CROSS JOIN LATERAL unnest(g.higher_taxon_keys) AS h(key)
)
INSERT INTO gbif_dependencies (import_id, key, from_key)
SELECT import_id,
    key,
    from_key
FROM expanded
WHERE NOT EXISTS (
        SELECT 1
        FROM taxa t
        WHERE t.gbif_id = key
    ) ON CONFLICT (import_id, key, from_key) DO NOTHING;

-- name: ListMissingGBIFKeys :many
SELECT DISTINCT d.key
FROM gbif_dependencies d
WHERE d.import_id = @import_id
    AND NOT EXISTS (
        SELECT 1
        FROM taxa t
        WHERE t.gbif_id = d.key
    )
    AND NOT EXISTS (
        SELECT 1
        FROM gbif_staging g
        WHERE g.key = d.key
    );

-- name: MaterializeTaxaFromGBIF :exec
-- This query materializes taxa from the GBIF staging table into the main taxa table, 
-- based on the dependencies identified for a given import batch. 
-- It ensures that only taxa that are not already present in the main taxa table are inserted.
--
-- ExpandGBIFDependencies should be run before this query to ensure that all necessary GBIF keys are identified.
-- The GBIF staging table should be populated with the relevant GBIF data before running this query.
-- Taxa to fetch can be identified with ListMissingGBIFKeys, which will return the GBIF keys 
-- that are needed for the import batch but are not yet present 
-- in the main taxa table or the GBIF staging table.
-- 
-- Fetched taxa can be inserted with InsertGBIFBatch, which will insert the taxa into the GBIF staging table,
-- if they are not already present.
INSERT INTO taxa (
        gbif_id,
        name,
        authorship,
        rank,
        status,
        parent_id,
        accepted_taxon_id
    )
SELECT g.key,
    g.canonical_name,
    g.authorship,
    g.rank::taxon_rank,
    g.status::taxon_status,
    parent.id as parent_id,
    accepted.id as accepted_taxon_id
FROM taxon_resolution r
    JOIN taxon_candidates c ON (
        r.resolved_candidate_id = c.id
        AND r.import_id = c.import_id
    )
    JOIN gbif_dependencies d ON (
        d.import_id = r.import_id
        AND d.from_key = c.gbif_id
    )
    JOIN gbif_staging g ON (g.key = d.key)
    LEFT JOIN taxa parent ON parent.gbif_id = g.parent_key
    LEFT JOIN taxa accepted ON accepted.gbif_id = g.accepted_key
WHERE r.import_id = @import_id
    AND g.rank = @rank
    AND (
        (
            @is_synonym = false
            AND g.accepted_key IS NULL
        )
        OR (
            @is_synonym = true
            AND g.accepted_key IS NOT NULL
        )
    ) ON CONFLICT (gbif_id) DO NOTHING;