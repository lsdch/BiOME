-- name: InitTaxonResolution :many
INSERT INTO taxon_resolution (import_id, input_name)
SELECT DISTINCT i.import_id,
    i.scientific_name
FROM import_samplings_occurrences i
WHERE i.import_id = @import_id
RETURNING *;

-- name: CleanUpTaxonResolution :exec
DELETE FROM taxon_resolution
WHERE import_id = @import_id;

-- name: GetTaxonResolution :many
SELECT *
FROM taxon_resolution
WHERE import_id = @import_id;

-- name: UpsertTaxonResolution :exec
INSERT INTO taxon_resolution (
        import_id,
        input_name,
        source,
        taxon_id,
        gbif_id,
        staging_id,
        status
    )
SELECT @import_id,
    @input_name,
    @source,
    @taxon_id,
    @gbif_id,
    @staging_id,
    @status ON CONFLICT (import_id, input_name) DO NOTHING;

-- name: ResolveTaxon :exec
UPDATE taxon_resolution
SET source = @source::taxon_match_source,
    match_type = @match_type::taxon_match_type,
    taxon_id = @taxon_id::UUID,
    gbif_id = @gbif_id::INTEGER,
    staging_id = @staging_id::UUID,
    status = @status::resolution_status
WHERE import_id = @import_id
    AND input_name = @input_name;

-- name: CreateCandidateTaxaNameExact :exec
-- Create candidate matches based on exact name matches, 
-- using either the scientific name or the canonical name.
INSERT INTO taxon_candidates (
        import_id,
        input_name,
        source,
        match_type,
        taxon_id,
        gbif_id,
        score,
        priority,
        name,
        authorship,
        rank,
        status
    )
SELECT i.import_id,
    i.input_name,
    'internal',
    'exact',
    t.id,
    t.gbif_id,
    1.0,
    100,
    t.name,
    t.authorship,
    t.rank,
    t.status
FROM taxon_resolution i
    JOIN taxa t ON i.input_name % t.scientific_name
    OR i.input_name % t.name
WHERE i.import_id = @import_id;

-- name: CreateCandidateTaxaFuzzy :exec
WITH candidates AS (
    SELECT i.import_id,
        i.input_name,
        t.id AS taxon_id,
        similarity(t.scientific_name, i.input_name) AS score,
        t.name,
        t.authorship,
        t.rank,
        t.status
    FROM taxon_resolution i
        JOIN taxa t ON t.scientific_name % i.input_name
    WHERE i.import_id = @import_id
)
INSERT INTO taxon_candidates (
        import_id,
        input_name,
        source,
        match_type,
        taxon_id,
        gbif_id,
        score,
        priority,
        name,
        authorship,
        rank,
        status
    )
SELECT c.import_id,
    c.input_name,
    'internal',
    'fuzzy',
    c.taxon_id,
    c.gbif_id,
    c.score,
    50,
    c.name,
    c.authorship,
    c.rank,
    c.status
FROM candidates c
WHERE c.score > COALESCE(NULLIF(@threshold::double precision, 0), 0.6);

-- name: MarkTaxaNeedingGBIFCandidates :exec
UPDATE taxon_resolution r
SET gbif_status = 'pending'
WHERE r.import_id = @import_id
    AND r.gbif_status NOT IN ('skipped', 'completed')
    AND NOT EXISTS (
        SELECT 1
        FROM taxon_candidates c
        WHERE c.import_id = r.import_id
            AND c.input_name = r.input_name
            AND c.source = 'internal'
            AND c.match_type = 'exact'
    )
    AND NOT EXISTS (
        SELECT 1
        FROM taxon_candidates c
        WHERE c.import_id = r.import_id
            AND c.input_name = r.input_name
            AND c.source = 'gbif'
    );

-- name: MarkTaxaGBIFImportCompleted :exec
UPDATE taxon_resolution r
SET gbif_status = 'completed'
WHERE r.import_id = @import_id
    AND r.gbif_status NOT IN ('skipped', 'completed')
    AND r.input_name = ANY(@input_names::TEXT []);

-- name: ListTaxaToFetchGBIFCandidates :many
SELECT DISTINCT i.taxon_name,
    i.taxon_scientific_name as full_input_name,
    i.taxon_rank
FROM import_samplings_occurrences i
    JOIN taxon_resolution r ON r.import_id = i.import_id
    AND r.input_name = i.taxon_scientific_name
WHERE i.import_id = @import_id
    AND r.gbif_status IN ('pending', 'failed')
ORDER BY i.taxon_scientific_name;

-- name: CleanupTaxonCandidates :exec
DELETE FROM taxon_candidates
WHERE import_id = @import_id;

-- name: ListAllTaxonCandidates :many
WITH staging AS (
    SELECT i.taxon_name,
        i.taxon_scientific_name,
        i.taxon_rank,
        i.taxon_authorship,
        ARRAY_AGG(i.row_number) AS row_numbers
    FROM import_samplings_occurrences i
    WHERE i.import_id = @import_id
    GROUP BY i.taxon_name,
        i.taxon_scientific_name,
        i.taxon_rank,
        i.taxon_authorship
),
candidates AS (
    SELECT s.taxon_name,
        s.taxon_authorship,
        s.taxon_rank,
        r.input_name,
        r.source,
        r.match_type,
        r.score,
        r.priority,
        COALESCE(t.name, g.canonical_name) AS resolved_name,
        COALESCE(t.authorship, g.authorship) AS resolved_authorship,
        COALESCE(t.rank, g.rank) AS resolved_rank,
        COALESCE(t.status, g.status) AS resolved_status,
        t.id AS resolved_taxon_id,
        r.gbif_id AS resolved_gbif_id
    FROM staging s
        JOIN taxon_candidates r ON r.input_name = s.taxon_scientific_name
        AND r.import_id = @import_id
        LEFT JOIN taxa t ON r.source = 'internal'
        AND t.id = r.taxon_id
        LEFT JOIN gbif_staging g ON r.source = 'gbif'
        AND g.key = r.gbif_id
)
SELECT *,
    ROW_NUMBER() OVER (
        PARTITION BY input_name
        ORDER BY priority DESC,
            score DESC
    ) AS rank
FROM candidates;

-- name: InsertTaxonCandidatesBatch :copyfrom
INSERT INTO taxon_candidates (
        import_id,
        input_name,
        source,
        match_type,
        taxon_id,
        gbif_id,
        score,
        priority,
        name,
        authorship,
        rank,
        status
    )
VALUES (
        @import_id,
        @input_name,
        @source,
        @match_type,
        @taxon_id,
        @gbif_id,
        @score,
        @priority,
        @name,
        @authorship,
        @rank,
        @status
    );

-- name: AutoResolveUnambiguousCandidates :exec
-- Automatically resolve candidates where there is a single best match above the priority threshold.
WITH ranked AS (
    SELECT c.*,
        ROW_NUMBER() OVER (
            PARTITION BY c.import_id,
            c.input_name
            ORDER BY c.priority DESC,
                c.score DESC
        ) AS rn,
        COUNT(*) FILTER (
            WHERE c.priority = (
                    MAX(c.priority) OVER (
                        PARTITION BY c.import_id,
                        c.input_name
                    )
                )
        ) OVER (
            PARTITION BY c.import_id,
            c.input_name
        ) AS best_count
    FROM taxon_candidates c
    WHERE c.import_id = @import_id
),
winners AS (
    SELECT *
    FROM ranked
    WHERE rn = 1
        AND best_count = 1
        AND priority >= 100
)
UPDATE taxon_resolution r
SET source = w.source,
    match_type = w.match_type,
    taxon_id = w.taxon_id,
    gbif_id = w.gbif_id,
    staging_id = w.staging_id,
    status = 'resolved'
FROM winners w
WHERE r.import_id = w.import_id
    AND r.input_name = w.input_name
    AND r.status = 'pending';