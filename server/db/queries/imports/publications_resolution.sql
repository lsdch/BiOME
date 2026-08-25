-- name: InitBibliographyResolution :batchone
WITH resolution AS (
    INSERT INTO publication_resolution (
            import_id,
            verbatim,
            doi,
            year,
            authors_raw,
            authors,
            title,
            journal
        )
    VALUES (
            @import_id,
            @verbatim,
            @doi,
            @year,
            @authors_raw,
            @authors,
            @title,
            @journal
        )
    RETURNING id
),
staging_occurrences AS (
    SELECT DISTINCT i.id AS occurrence_id,
        i.row_number
    FROM import_samplings_occurrences i
    WHERE i.import_id = @import_id
        AND i.row_number = ANY(@row_numbers::int [])
),
resolution_links AS (
    INSERT INTO occurrences_staging_publications (import_id, occurrence_id, resolution_id)
    SELECT @import_id,
        o.occurrence_id,
        r.id
    FROM resolution r
        CROSS JOIN staging_occurrences o
    RETURNING *
),
missing_rows AS (
    SELECT rn
    FROM unnest(@row_numbers::int []) AS rn
    WHERE NOT EXISTS (
            SELECT 1
            FROM staging_occurrences o
            WHERE o.row_number = rn
        )
)
SELECT COALESCE(array_agg(rn), '{}'::int [])::int [] AS missing_rows
FROM missing_rows;


-- name: ResolvePublication :exec 
UPDATE publication_resolution
SET resolved_candidate_id = @resolved_candidate_id::uuid,
    status = 'user_resolved'
WHERE import_id = @import_id
    AND id = @resolution_id;



-- name: GenerateBibliographyManualCandidates :exec
WITH resolution AS (
    SELECT *
    FROM publication_resolution r
    WHERE r.import_id = @import_id
        AND r.doi IS NULL
        AND NOT EXISTS (
            SELECT 1
            FROM publication_candidates c
                JOIN publications p ON p.id = c.internal_id
            WHERE p.verbatim = r.verbatim
        )
),
staging AS (
    INSERT INTO publications_staging (
            origin_resolution_id,
            doi,
            verbatim,
            authors,
            year,
            title,
            journal,
            source
        )
    SELECT r.id,
        NULL::doi,
        r.verbatim,
        r.authors,
        r.year,
        r.title,
        r.journal,
        'manual'::publication_source
    FROM resolution r
    RETURNING *
)
INSERT INTO publication_candidates (
        import_id,
        resolution_id,
        staging_id,
        score,
        match_type,
        source
    )
SELECT r.import_id,
    r.id,
    s.id,
    -1.0,
    'verbatim'::pub_match_type,
    'manual'::publication_candidate_source
FROM resolution r
    JOIN staging s ON (s.origin_resolution_id = r.id) ON CONFLICT DO NOTHING;


-- name: GenerateBibliographyInternalCandidates :exec
WITH resolution AS (
    SELECT id,
        doi,
        verbatim
    FROM publication_resolution
    WHERE import_id = @import_id
        AND status = 'pending'
),
internal_candidates AS (
    SELECT r.id AS resolution_id,
        p.id AS internal_id,
        CASE
            WHEN r.doi IS NOT NULL
            AND r.doi = p.doi THEN 'doi'::pub_match_type
            ELSE 'verbatim'::pub_match_type
        END AS match_type,
        CASE
            WHEN r.doi IS NOT NULL
            AND r.doi = p.doi THEN 100
            ELSE similarity(p.verbatim, r.verbatim) * 100
        END AS score
    FROM resolution r
        JOIN publications p ON (
            (
                r.doi IS NOT NULL
                AND p.doi = r.doi
            )
            OR (
                r.verbatim IS NOT NULL
                AND similarity(p.verbatim, r.verbatim) > 0.5
            )
        )
)
INSERT INTO publication_candidates (
        import_id,
        resolution_id,
        internal_id,
        score,
        match_type,
        source
    )
SELECT @import_id,
    ic.resolution_id,
    ic.internal_id,
    ic.score,
    ic.match_type,
    'internal'::publication_candidate_source
FROM internal_candidates ic ON CONFLICT DO NOTHING;

-- name: ResolvePendingToManualCandidates :exec
WITH resolution AS (
    SELECT *
    FROM publication_resolution r
    WHERE r.import_id = @import_id
        AND r.status = 'pending'
),
manual_candidates AS (
    SELECT *
    FROM publication_candidates pc
    WHERE pc.import_id = @import_id
        AND pc.source = 'manual'::publication_candidate_source
)
UPDATE publication_resolution r
SET status = 'user_resolved',
    resolved_candidate_id = mc.id
FROM manual_candidates mc
WHERE r.id = mc.resolution_id;

-- name: AutoResolveBibliography :exec
WITH resolution AS (
    SELECT id
    FROM publication_resolution r
    WHERE r.import_id = @import_id
        AND r.status = 'pending'
),
ranked_candidates AS (
    SELECT pc.id,
        pc.resolution_id,
        pc.match_type,
        pc.score,
        ROW_NUMBER() OVER (
            PARTITION BY pc.resolution_id
            ORDER BY pc.score DESC,
                pc.id
        ) AS rn,
        LEAD(pc.score) OVER (
            PARTITION BY pc.resolution_id
            ORDER BY pc.score DESC,
                pc.id
        ) AS next_score
    FROM publication_candidates pc
        JOIN resolution r ON pc.resolution_id = r.id
),
resolved_candidates AS (
    SELECT id,
        resolution_id
    FROM ranked_candidates
    WHERE rn = 1
        AND (
            match_type = 'doi'::pub_match_type
            OR (
                score >= @score_threshold
                AND COALESCE(score - next_score, score) >= @score_margin
            )
        )
)
UPDATE publication_resolution r
SET status = 'auto_resolved',
    resolved_candidate_id = rc.id
FROM resolved_candidates rc
WHERE r.id = rc.resolution_id;

-- name: ResolveBibliographyManualCandidates :exec
WITH single_candidates AS (
    SELECT resolution_id
    FROM publication_candidates
    WHERE import_id = @import_id
    GROUP BY resolution_id
    HAVING COUNT(*) = 1
),
manual_candidates AS (
    SELECT pc.resolution_id,
        pc.id
    FROM publication_candidates pc
        JOIN single_candidates sc ON sc.resolution_id = pc.resolution_id
    WHERE pc.import_id = @import_id
        AND pc.source = 'manual'::publication_candidate_source
)
UPDATE publication_resolution r
SET status = 'user_resolved',
    resolved_candidate_id = mc.id
FROM manual_candidates mc
WHERE r.id = mc.resolution_id
    AND r.import_id = @import_id
    AND r.status = 'pending';

  
-- name: ListDOIsToFetch :many
SELECT DISTINCT r.doi::doi
FROM publication_resolution r
WHERE r.import_id = @import_id
    AND r.status = 'pending'
    AND r.doi IS NOT NULL
    AND NOT EXISTS (
        SELECT 1
        FROM publication_candidates pc
        WHERE pc.resolution_id = r.id
            AND pc.match_type = 'doi'::pub_match_type
    )
    AND NOT EXISTS (
        SELECT 1
        FROM publications p
        WHERE p.doi = r.doi
    )
    AND NOT EXISTS (
        SELECT 1
        FROM publications_staging s
        WHERE s.doi = r.doi
    );

-- name: ListVerbatimToFetch :many
SELECT DISTINCT r.verbatim::text
FROM publication_resolution r
WHERE r.import_id = @import_id
    AND r.status = 'pending'
    AND r.verbatim IS NOT NULL
    AND r.doi IS NULL
    AND NOT EXISTS (
        SELECT 1
        FROM publication_candidates pc
        WHERE pc.resolution_id = r.id
            AND pc.match_type = 'verbatim'::pub_match_type
            AND pc.source = 'crossref'::publication_candidate_source
    );

-- name: StagePublications :copyfrom
INSERT INTO publications_staging (
        doi,
        verbatim,
        authors,
        year,
        title,
        journal,
        source
    )
VALUES (
        @doi,
        @verbatim,
        @authors,
        @year,
        @title,
        @journal,
        @source
    );

    
-- name: GenerateBibliographyExternalCandidates :exec
WITH resolution AS (
    SELECT id,
        import_id,
        doi,
        verbatim
    FROM publication_resolution r
    WHERE r.import_id = @import_id
        AND r.status = 'pending'
),
doi_matches AS (
    SELECT r.id AS resolution_id,
        r.import_id,
        s.id AS staging_id,
        s.source::text AS source,
        'doi'::pub_match_type AS match_type,
        100.0 AS score
    FROM resolution r
        JOIN publications_staging s ON s.doi = r.doi
    WHERE r.doi IS NOT NULL
) -- verbatim_matches AS (
--     SELECT r.id AS resolution_id,
--         r.import_id,
--         s.id AS staging_id,
--         s.source::text AS source,
--         'verbatim'::pub_match_type AS match_type,
--         sim.score
--     FROM resolution r
--         JOIN publications_staging s ON (
--             r.verbatim IS NOT NULL
--             AND s.source = 'crossref'::publication_source
--         )
--         CROSS JOIN LATERAL (
--             SELECT similarity(s.verbatim, r.verbatim) * 100 AS score
--         ) sim
--     WHERE sim.score > 30
-- ),
-- staging_candidates AS (
--     SELECT *
--     FROM doi_matches
--     UNION ALL
--     SELECT vm.*
--     FROM verbatim_matches vm
--     WHERE NOT EXISTS (
--             SELECT 1
--             FROM doi_matches dm
--             WHERE dm.resolution_id = vm.resolution_id
--                 AND dm.staging_id = vm.staging_id
--         )ff55a73e-02f6-4364-b31f-77f29525305b
-- )
INSERT INTO publication_candidates (
        import_id,
        resolution_id,
        staging_id,
        score,
        match_type,
        source
    )
SELECT import_id,
    resolution_id,
    staging_id,
    score,
    match_type,
    source::publication_candidate_source
FROM doi_matches ON CONFLICT DO NOTHING;

-- name: ListPublicationResolutions :many
SELECT r.*
FROM publication_resolution r
WHERE r.import_id = @import_id
ORDER BY r.id;

-- name: ListPublicationCandidates :many
SELECT sqlc.embed(pc),
    COALESCE(p.doi, s.doi) AS doi,
    COALESCE(p.verbatim, s.verbatim) AS verbatim,
    COALESCE(p.authors, s.authors) AS authors,
    COALESCE(p.year, s.year) AS year,
    COALESCE(p.title, s.title) AS title,
    COALESCE(p.journal, s.journal) AS journal
FROM publication_candidates pc
    JOIN publication_resolution r ON pc.resolution_id = r.id
    LEFT JOIN publications p ON (
        pc.source = 'internal'
        AND pc.internal_id = p.id
    )
    LEFT JOIN publications_staging s ON (
        pc.source IN ('crossref', 'manual')
        AND pc.staging_id = s.id
    )
WHERE r.import_id = @import_id;

-- name: MaterializeBibliography :exec
WITH materialized AS (
    INSERT INTO publications (
            id,
            authors,
            year,
            title,
            journal,
            verbatim,
            doi
        )
    SELECT s.id,
        s.authors,
        s.year,
        s.title,
        s.journal,
        s.verbatim,
        s.doi
    FROM publications_staging s
        JOIN publication_candidates pc ON pc.staging_id = s.id
    WHERE pc.import_id = @import_id ON CONFLICT DO NOTHING
    RETURNING *
)
INSERT INTO occurrences_publications (occurrence_id, publication_id)
SELECT osp.occurrence_id,
    COALESCE(pc.internal_id, m.id)
FROM occurrences_staging_publications osp
    JOIN publication_resolution r ON r.id = osp.resolution_id
    JOIN publication_candidates pc ON pc.id = r.resolved_candidate_id
    LEFT JOIN materialized m ON m.id = pc.staging_id
WHERE osp.import_id = @import_id;