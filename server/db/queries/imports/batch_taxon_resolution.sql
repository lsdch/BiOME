-- name: CheckTaxaConsistencyInImport :many
SELECT input_name,
    count(*),
    array_agg(DISTINCT taxon_authorship)::text [] AS authorships,
    array_agg(DISTINCT taxon_rank)::text [] AS ranks
FROM (
        SELECT DISTINCT i.import_id,
            i.taxon_name::citext AS input_name,
            i.taxon_authorship,
            i.taxon_rank
        FROM import_samplings_occurrences i
        WHERE i.import_id = $1
    ) t
GROUP BY input_name
HAVING count(*) > 1;

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

-- name: InitSamplingTargetResolution :exec
WITH all_sampling_targets AS (
    SELECT DISTINCT import_id,
        sampling_hash,
        sampling_target
    FROM import_samplings_occurrences iso
        CROSS JOIN LATERAL unnest(sampling_targets) AS sampling_target
    WHERE iso.import_id = @import_id
        AND sampling_target <> ''
),
new_resolutions AS (
    INSERT INTO taxon_resolution (import_id, input_name, sampling_target)
    SELECT import_id,
        sampling_target,
        true
    FROM all_sampling_targets ON CONFLICT (import_id, input_name) DO NOTHING
    RETURNING *
)
INSERT INTO sampling_target_resolution (
        import_id,
        sampling_hash,
        resolution_id
    )
SELECT st.import_id,
    st.sampling_hash,
    r.id
FROM all_sampling_targets st
    JOIN taxon_resolution r ON r.import_id = st.import_id
    AND r.input_name = st.sampling_target ON CONFLICT (import_id, sampling_hash, resolution_id) DO NOTHING;

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
SELECT sqlc.embed(t),
    parent.input_name as from_resolution_name
FROM taxon_resolution t
    LEFT JOIN taxon_resolution parent ON t.from_resolution_id = parent.id
WHERE t.import_id = @import_id;

-- name: UpsertTaxonResolution :exec
INSERT INTO taxon_resolution (
        import_id,
        input_name,
        input_authorship,
        input_rank,
        resolved_candidate_id,
        status
    )
VALUES (
        @import_id,
        @input_name,
        @input_authorship,
        @input_rank,
        @resolved_candidate_id,
        @status
    ) ON CONFLICT (import_id, input_name) DO NOTHING;

-- name: ResolveTaxon :exec
UPDATE taxon_resolution
SET resolved_candidate_id = @candidate_id::uuid,
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
    AND r.gbif_status = 'pending'
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
SET gbif_status = CASE
        WHEN EXISTS (
            SELECT 1
            FROM taxon_candidates c
            WHERE c.import_id = r.import_id
                AND c.resolution_id = r.id
        ) THEN 'completed'::taxon_gbif_status
        ELSE 'no_candidates'::taxon_gbif_status
    END
WHERE r.import_id = @import_id
    AND r.gbif_status = 'pending';

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
        COALESCE(t.name, g.canonical_name, s.name) AS resolved_name,
        COALESCE(t.authorship, g.authorship, s.authorship) AS resolved_authorship,
        COALESCE(t.rank, g.rank::taxon_rank, s.rank) AS resolved_rank,
        COALESCE(t.status, g.status::taxon_status, s.status) AS resolved_status
    FROM taxon_candidates c
        LEFT JOIN taxa t ON (
            c.source = 'internal'
            AND t.id = c.taxon_id
        )
        LEFT JOIN gbif_staging g ON (
            c.source = 'gbif'
            AND g.key = c.gbif_id
        )
        LEFT JOIN taxa_staging s ON (
            c.source = 'manual'
            AND s.id = c.staging_id
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
        AND c.priority >= @threshold
    GROUP BY c.import_id,
        c.resolution_id
    HAVING COUNT(*) = 1
)
UPDATE taxon_resolution r
SET resolved_candidate_id = w.id,
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

-- name: SetNeedsResolutionForUnresolvedCandidates :exec
-- Mark taxon resolutions as needing user resolution if they have candidates but are still unresolved.
UPDATE taxon_resolution r
SET status = 'needs_decision'
WHERE r.import_id = @import_id
    AND r.status = 'pending';

-- name: UpdateMaterializedGBIFCandidates :exec
-- Update the taxon_id field in taxon_candidates for candidates that have a matching gbif_id in the taxa table.
-- This is necessary because the taxa table is populated after the taxon_candidates table, 
-- and we need to link the candidates to the actual taxa records.
UPDATE taxon_candidates c
SET taxon_id = t.id
FROM taxa t
WHERE c.gbif_id = t.gbif_id
    AND c.import_id = @import_id
    AND c.taxon_id IS NULL;



-- name: ListTaxonResolutionsWithoutCandidates :many
SELECT r.*
FROM taxon_resolution r
WHERE r.import_id = @import_id
    AND r.status = 'pending'
    AND NOT EXISTS (
        SELECT 1
        FROM taxon_candidates c
        WHERE c.import_id = r.import_id
            AND c.resolution_id = r.id
    );

-- name: InsertTaxaStaging :batchexec
-- Insert a new record into the taxa_staging table for a given taxon resolution,
-- and create a corresponding taxon_candidates record for it.
-- This is used when a user manually adds a new taxon that is not found in the existing taxa or GBIF data.
-- Resolution for the parent taxon is also created if it does not already exist, 
-- or linked to the existing resolution if it does, based on the provided parent name and rank.
-- AutoResolveUnambiguousCandidates should be run after this to resolve any unambiguous candidates that may have been created.
WITH inserted_parent_resolution AS (
    INSERT INTO taxon_resolution (
            import_id,
            input_name,
            input_authorship,
            input_rank,
            status,
            from_resolution_id
        )
    VALUES (
            @import_id,
            @parent_name,
            NULL,
            @parent_rank,
            'pending',
            @resolution_id
        ) ON CONFLICT (import_id, input_name) DO NOTHING
    RETURNING id
),
staged_taxon AS (
    INSERT INTO taxa_staging (
            import_id,
            name,
            authorship,
            rank,
            status,
            parent_resolution_id
        )
    VALUES (
            @import_id,
            @name,
            @authorship,
            @taxon_rank,
            @taxon_status,
            COALESCE(
                (
                    SELECT id
                    FROM inserted_parent_resolution
                ),
                (
                    SELECT parent_res.id
                    FROM taxon_resolution parent_res
                    WHERE import_id = @import_id
                        AND input_name = @parent_name
                    LIMIT 1
                )
            )
        )
    RETURNING *
)
INSERT INTO taxon_candidates (
        import_id,
        resolution_id,
        source,
        match_type,
        staging_id,
        priority,
        name,
        authorship,
        rank,
        status
    )
SELECT @import_id,
    @resolution_id,
    'manual',
    'exact',
    s.id,
    100,
    s.name,
    s.authorship,
    s.rank,
    s.status
FROM staged_taxon s;

-- name: MaterializeTaxaStaging :exec
-- Materialize taxa from the taxa_staging table into the main taxa table,
-- based on the resolutions identified for a given import batch.
-- Materialization of GBIF taxa must have been completed before this query is run, 
-- as it relies on the resolved candidates from the GBIF data.
-- This query must be run for each rank separately, 
-- as the parent taxa must be materialized before the child taxa can be materialized.
-- SyncMaterializedTaxa must be run after each rank is materialized 
-- to update the taxon_candidates table with the newly created taxa IDs.
WITH resolved_candidates AS (
    SELECT s.*,
        parent_candidate.taxon_id AS parent_id
    FROM taxon_candidates c
        JOIN taxon_resolution r ON (
            r.import_id = c.import_id
            AND r.resolved_candidate_id = c.id
        )
        JOIN taxa_staging s ON (c.staging_id = s.id)
        JOIN taxon_resolution parent_resolution ON (
            parent_resolution.id = s.parent_resolution_id
            AND parent_resolution.import_id = s.import_id
        )
        JOIN taxon_candidates parent_candidate ON (
            parent_candidate.id = parent_resolution.resolved_candidate_id
            AND parent_candidate.import_id = s.import_id
        )
    WHERE c.import_id = @import_id
        AND c.source = 'manual'
        AND c.rank = @rank
)
INSERT INTO taxa (
        name,
        authorship,
        rank,
        status,
        parent_id
    )
SELECT name,
    authorship,
    rank,
    status,
    parent_id
FROM resolved_candidates ON CONFLICT (name, COALESCE(authorship, '')) DO NOTHING;


-- name: SyncMaterializedTaxa :exec
UPDATE taxon_candidates c
SET taxon_id = t.id
FROM taxa_staging s
    JOIN taxa t ON (
        t.name = s.name
        AND t.authorship IS NOT DISTINCT
        FROM s.authorship
            AND t.rank = s.rank
    )
WHERE c.import_id = @import_id
    AND c.staging_id = s.id
    AND c.rank = @rank
    AND c.source = 'manual';