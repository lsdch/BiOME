-- name: InsertGBIFBatch :copyfrom
INSERT INTO gbif_staging (
        input_name,
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
        @input_name::TEXT,
        @key::INTEGER,
        @parent::TEXT,
        @parent_key::INTEGER,
        @canonical_name::TEXT,
        @scientific_name::TEXT,
        @status::TEXT,
        @rank::TEXT,
        @name_type::TEXT,
        @kingdom_key::INTEGER,
        @phylum_key::INTEGER,
        @class_key::INTEGER,
        @order_key::INTEGER,
        @family_key::INTEGER,
        @genus_key::INTEGER,
        @species_key::INTEGER,
        @higher_taxon_keys::INTEGER [],
        @higher_taxon_names::TEXT [],
        @authorship::TEXT,
        @num_descendants::INTEGER,
        @accepted_key::INTEGER,
        @accepted_name::TEXT
    );

-- name: CleanUpGBIFDependencies :exec
DELETE FROM gbif_dependencies d
WHERE d.import_id = @import_id;

-- name: ExpandGBIFDependencies :exec
WITH resolution AS (
    SELECT DISTINCT r.import_id,
        r.gbif_id
    FROM taxon_resolution r
        JOIN gbif_staging g ON g.key = taxon_resolution.gbif_id
    WHERE r.import_id = @import_id
        AND gbif_id IS NOT NULL
),
expanded AS (
    -- 1. direct gbif_id
    SELECT r.import_id,
        r.gbif_id AS key
    FROM resolution r
    UNION ALL
    -- 2. accepted key
    SELECT r.import_id,
        g.accepted_key AS key
    FROM resolution r
        JOIN gbif_staging g ON g.key = r.gbif_id
    WHERE g.accepted_key IS NOT NULL
    UNION ALL
    -- 3. higher taxonomy expansion
    SELECT r.import_id,
        h.key
    FROM resolution r
        JOIN gbif_staging g ON g.key = r.gbif_id
        CROSS JOIN LATERAL unnest(g.higher_taxon_keys) AS h(key)
),
deduplicated AS (
    SELECT DISTINCT import_id,
        key
    FROM expanded
    WHERE key IS NOT NULL
)
INSERT INTO gbif_dependencies (import_id, key)
SELECT import_id,
    key
FROM deduplicated
WHERE NOT EXISTS (
        SELECT 1
        FROM taxa t
        WHERE t.gbif_key = key
    ) ON CONFLICT (import_id, key) DO NOTHING;

-- name: ListMissingGBIFKeys :many
SELECT d.key
FROM gbif_dependencies d
WHERE d.import_id = @import_id
    AND NOT EXISTS (
        SELECT 1
        FROM taxa t
        WHERE t.gbif_key = d.key
    )
    AND NOT EXISTS (
        SELECT 1
        FROM gbif_staging g
        WHERE g.key = d.key
    );

-- name: InsertTaxaFromGBIF :exec
INSERT INTO taxa (
        gbif_id,
        name,
        authorship,
        rank,
        status,
        parent_id,
        accepted_id
    )
SELECT g.key,
    g.canonical_name,
    g.authorship,
    g.rank,
    g.status,
    parent.id,
    accepted.id
FROM gbif_staging g
    JOIN gbif_dependencies d ON d.key = g.key
    LEFT JOIN taxa parent ON parent.gbif_key = g.parent_key
    LEFT JOIN taxa accepted ON accepted.gbif_key = g.accepted_key
WHERE d.import_id = @import_id
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

-- name: InsertTaxonStaging :exec
INSERT INTO taxa_staging (
        import_id,
        name,
        authorship,
        rank,
        status,
        parent_source,
        parent_taxa_id,
        parent_gbif_id,
        parent_input_name
    )
SELECT @import_id,
    @name,
    @authorship,
    @rank,
    @status,
    @parent_source,
    @parent_taxa_id,
    @parent_gbif_id,
    @parent_input_name;

-- name: MaterializeTaxaStaging :exec
-- INSERT INTO taxa (
--         name,
--         authorship,
--         rank,
--         status,
--         parent_id
--     )
-- SELECT s.name,
--     s.authorship,
--     s.rank,
--     s.status,
--     parent.id
-- FROM taxa_staging s
--     JOIN taxon_resolution r ON r.staging_id = s.id
--     LEFT JOIN taxa parent ON parent.id = s.parent_taxa_id
--     OR parent.gbif_id = s.parent_gbif_id
-- WHERE r.import_id = @import_id
--     AND r.source = 'manual'
--     AND s.rank = @rank