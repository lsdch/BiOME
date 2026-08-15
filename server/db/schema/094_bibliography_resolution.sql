CREATE TYPE publication_source AS ENUM ('crossref', 'manual');

CREATE TYPE publication_candidate_source AS ENUM ('internal', 'crossref', 'manual');

CREATE TABLE publications_staging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doi DOI,
    verbatim TEXT NOT NULL,
    authors TEXT [],
    year INTEGER,
    title TEXT,
    journal TEXT,
    source publication_source NOT NULL
);

CREATE TYPE pub_match_type AS ENUM ('doi', 'verbatim');

-- A repository of candidate publications that imported publications can be resolved to. 
-- This table is populated from the crossref_staging table and the publications table.
CREATE TABLE publication_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- ==========================
    -- INGESTION CONTEXT
    -- ==========================
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    -- ADDED AFTER PUBLICATION RESOLUTION TABLE IS CREATED : 
    -- resolution_id UUID NOT NULL REFERENCES publication_resolution (id) ON DELETE CASCADE,
    -- ==========================
    -- CANDIDATE METADATA
    -- ==========================
    match_type pub_match_type NOT NULL,
    internal_id UUID REFERENCES publications (id) ON DELETE CASCADE,
    staging_id UUID REFERENCES publications_staging (id) ON DELETE CASCADE,
    score REAL NOT NULL,
    source publication_candidate_source NOT NULL,
    CONSTRAINT internal_or_staging_check CHECK (
        (
            internal_id IS NOT NULL
            OR staging_id IS NOT NULL
        )
    )
);

CREATE TABLE publication_resolution (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    status resolution_status NOT NULL DEFAULT 'pending',
    resolved_candidate_id UUID REFERENCES publication_candidates (id) ON DELETE CASCADE,
    doi DOI,
    verbatim TEXT,
    authors TEXT [],
    authors_raw TEXT,
    year INTEGER,
    title TEXT,
    journal TEXT,
    CONSTRAINT doi_or_verbatim_check CHECK (
        (
            doi IS NOT NULL
            OR verbatim IS NOT NULL
        )
    )
);

ALTER TABLE publication_candidates
ADD COLUMN resolution_id UUID NOT NULL REFERENCES publication_resolution (id) ON DELETE CASCADE,
    ADD CONSTRAINT candidate_unique UNIQUE (resolution_id, internal_id, staging_id);

CREATE UNIQUE INDEX candidate_unique_manual ON publication_candidates (resolution_id, staging_id)
WHERE source = 'manual';

ALTER TABLE publications_staging
ADD COLUMN origin_resolution_id UUID REFERENCES publication_resolution (id) ON DELETE CASCADE;


-- A many-to-many association table linking occurrences to the publications resolution table
CREATE TABLE occurrences_staging_publications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- =========================
    -- INGESTION CONTEXT
    -- =========================
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    occurrence_id ULID NOT NULL REFERENCES import_samplings_occurrences (id) ON DELETE CASCADE,
    resolution_id UUID REFERENCES publication_resolution (id) ON DELETE CASCADE,
    CONSTRAINT occurrences_staging_publications_unique UNIQUE (import_id, occurrence_id, resolution_id)
);