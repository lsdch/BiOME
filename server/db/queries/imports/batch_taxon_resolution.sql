-- name: InitTaxonResolution :many
INSERT INTO taxon_resolution (
        import_id,
        input_name,
        input_authorship,
        input_rank
    )
SELECT DISTINCT i.import_id,
    i.taxon_name,
    i.taxon_authorship,
    i.taxon_rank
FROM import_samplings_occurrences i
WHERE i.import_id = @import_id
RETURNING *;

-- name: LinkTaxonResolutions :exec
UPDATE import_samplings_occurrences i
SET taxon_resolution_id = r.id
FROM taxon_resolution r
WHERE i.import_id = @import_id
    AND i.import_id = r.import_id
    AND i.taxon_name = r.input_name
    AND (
        i.taxon_authorship = r.input_authorship
        OR (
            i.taxon_authorship IS NULL
            AND r.input_authorship IS NULL
        )
    );

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
        input_authorship,
        input_rank,
        resolved_to,
        status
    )
VALUES (
        @import_id,
        @input_name,
        @input_authorship,
        @input_rank,
        @resolved_to,
        @status
    ) ON CONFLICT (import_id, input_name) DO NOTHING;

-- name: ResolveTaxon :exec
UPDATE taxon_resolution
SET resolved_to = @resolved_to::uuid,
    status = 'user_resolved'
WHERE import_id = @import_id
    AND id = @resolution_id;

-- name: CreateCandidateTaxaNameExact :exec
-- Create candidate matches based on exact name matches, 
-- using either the scientific name or the canonical name.
INSERT INTO taxon_candidates (
        import_id,
        resolution_id,
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
SELECT DISTINCT i.import_id,
    i.id,
    'internal'::taxon_match_source,
    'exact'::taxon_match_type,
    t.id,
    t.gbif_id,
    1.0,
    100,
    t.name,
    t.authorship,
    t.rank,
    t.status
FROM taxon_resolution i
    JOIN taxa t ON (
        i.scientific_name = t.scientific_name
        OR (
            i.input_authorship IS NULL
            AND i.input_name = t.name
        )
    )
WHERE i.import_id = @import_id;

-- name: CreateCandidateTaxaFuzzy :exec
WITH candidates AS (
    SELECT r.import_id,
        r.id AS resolution_id,
        t.id AS taxon_id,
        t.gbif_id,
        similarity(t.scientific_name, r.scientific_name) AS score,
        t.name,
        t.authorship,
        t.rank,
        t.status
    FROM taxon_resolution r
        JOIN taxa t ON t.scientific_name % r.scientific_name
    WHERE r.import_id = @import_id
)
INSERT INTO taxon_candidates (
        import_id,
        resolution_id,
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
    c.resolution_id,
    'internal'::taxon_match_source,
    'fuzzy'::taxon_match_type,
    c.taxon_id,
    c.gbif_id,
    c.score,
    50,
    c.name,
    c.authorship,
    c.rank,
    c.status
FROM candidates c
WHERE c.score > COALESCE(NULLIF(@threshold::double precision, 0), 0.6) ON CONFLICT DO NOTHING;

-- name: MarkTaxaNeedingGBIFCandidates :exec
UPDATE taxon_resolution r
SET gbif_status = 'pending'
WHERE r.import_id = @import_id
    AND r.gbif_status NOT IN ('skipped', 'completed')
    AND NOT EXISTS (
        SELECT 1
        FROM taxon_candidates c
        WHERE c.import_id = r.import_id
            AND c.resolution_id = r.id
            AND (
                (
                    c.source = 'internal'
                    AND c.match_type = 'exact'
                )
                OR c.source = 'gbif'
            )
    );

-- name: MarkTaxaGBIFImportCompleted :exec
UPDATE taxon_resolution r
SET gbif_status = 'completed'
WHERE r.import_id = @import_id
    AND r.gbif_status NOT IN ('skipped', 'completed')
    AND r.id = ANY(@resolutions_ids::UUID []);

-- name: ListResolutionsToFetchGBIFCandidates :many
SELECT r.*
FROM taxon_resolution r
WHERE r.import_id = @import_id
    AND r.gbif_status IN ('pending', 'failed')
ORDER BY r.scientific_name;

-- name: CleanupTaxonCandidates :exec
DELETE FROM taxon_candidates
WHERE import_id = @import_id;

-- name: ListAllTaxonCandidates :many
WITH candidates AS (
    SELECT c.id,
        c.resolution_id,
        c.source,
        c.match_type,
        c.score,
        c.priority,
        t.id AS resolved_taxon_id,
        c.gbif_id AS resolved_gbif_id,
        COALESCE(t.name, g.canonical_name) AS resolved_name,
        COALESCE(t.authorship, g.authorship) AS resolved_authorship,
        COALESCE(t.rank, g.rank::taxon_rank) AS resolved_rank,
        COALESCE(t.status, g.status::taxon_status) AS resolved_status
    FROM taxon_candidates c
        LEFT JOIN taxa t ON (
            c.source = 'internal'
            AND t.id = c.taxon_id
        )
        LEFT JOIN gbif_staging g ON (
            c.source = 'gbif'
            AND g.key = c.gbif_id
        )
    WHERE c.import_id = @import_id
)
SELECT *,
    ROW_NUMBER() OVER (
        PARTITION BY resolution_id
        ORDER BY priority DESC,
            score DESC
    ) AS rank
FROM candidates;

-- name: InsertTaxonCandidatesBatch :batchexec
INSERT INTO taxon_candidates (
        import_id,
        resolution_id,
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
        @resolution_id,
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
    ) ON CONFLICT DO NOTHING;

-- name: AutoResolveUnambiguousCandidates :exec
-- Automatically resolve candidates where there is a single best match above the priority threshold.
WITH winners AS (
    SELECT c.import_id,
        c.resolution_id,
        (array_agg(c.id)) [1] AS id
    FROM taxon_candidates c
    WHERE c.import_id = @import_id
        AND c.priority >= 100
    GROUP BY c.import_id,
        c.resolution_id
    HAVING COUNT(*) = 1
)
UPDATE taxon_resolution r
SET resolved_to = w.id,
    status = 'auto_resolved',
    gbif_status = (
        CASE
            WHEN r.gbif_status = 'completed' THEN r.gbif_status
            ELSE 'skipped'
        END
    )
FROM winners w
WHERE r.import_id = w.import_id
    AND r.id = w.resolution_id
    AND r.status = 'pending';

-- name: UpdateMaterializedTaxonCandidates :exec
-- Update the taxon_id field in taxon_candidates for candidates that have a matching gbif_id in the taxa table.
-- This is necessary because the taxa table is populated after the taxon_candidates table, 
-- and we need to link the candidates to the actual taxa records.
UPDATE taxon_candidates c
SET taxon_id = t.id
FROM taxa t
WHERE c.gbif_id = t.gbif_id
    AND c.import_id = @import_id
    AND c.taxon_id IS NULL;