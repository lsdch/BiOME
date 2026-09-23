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

CREATE INDEX publications_staging_doi_idx ON publications_staging (doi);

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
    resolved_candidate_id UUID,
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
    ADD CONSTRAINT publication_candidates_resolution_id_unique UNIQUE (resolution_id, id);

ALTER TABLE publication_resolution
ADD CONSTRAINT publication_resolution_resolved_candidate_fk FOREIGN KEY (id, resolved_candidate_id) REFERENCES publication_candidates (resolution_id, id) ON DELETE
SET NULL (resolved_candidate_id),
    ADD CONSTRAINT publication_resolution_resolved_status_check CHECK (
        status NOT IN ('auto_resolved', 'user_resolved')
        OR resolved_candidate_id IS NOT NULL
    );

-- A deleted or cleared selection must become eligible for resolution again.
CREATE FUNCTION publication_resolution_reset_status() RETURNS trigger AS $$ BEGIN IF NEW.resolved_candidate_id IS NULL
AND NEW.status IN ('auto_resolved', 'user_resolved') THEN NEW.status := 'pending';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER publication_resolution_reset_status_trigger BEFORE
UPDATE OF resolved_candidate_id ON publication_resolution FOR EACH ROW EXECUTE FUNCTION publication_resolution_reset_status();


CREATE UNIQUE INDEX publication_candidates_internal_uq ON publication_candidates (resolution_id, internal_id)
WHERE internal_id IS NOT NULL;

CREATE UNIQUE INDEX publication_candidates_staging_uq ON publication_candidates (resolution_id, staging_id)
WHERE staging_id IS NOT NULL;

ALTER TABLE publications_staging
ADD COLUMN origin_resolution_id UUID REFERENCES publication_resolution (id) ON DELETE CASCADE;

-- Prevent multiple manual candidates from being created for the same publication resolution.
CREATE UNIQUE INDEX publications_staging_manual_origin_uq ON publications_staging (origin_resolution_id)
WHERE source = 'manual';

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