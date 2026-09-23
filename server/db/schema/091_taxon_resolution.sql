CREATE TABLE IF NOT EXISTS gbif_dependencies (
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    key INTEGER NOT NULL,
    from_key INTEGER NOT NULL REFERENCES gbif_staging (key) ON DELETE CASCADE,
    PRIMARY KEY (import_id, key, from_key)
);


CREATE TYPE resolution_status AS ENUM (
    -- resolution has yet to be processed, e.g. candidates fetched and staged
    'pending',
    -- resolution has been processed and a natural candidate has been selected
    'auto_resolved',
    -- resolution has been processed and a candidate was selected by the user
    'user_resolved',
    -- resolution is ambiguous and requires user intervention to select a candidate
    'needs_decision'
);

CREATE TYPE taxon_gbif_status AS ENUM (
    -- skipped because the taxon was resolved internally or manually
    'skipped',
    -- waiting for GBIF candidates to be fetched
    'pending',
    -- GBIF candidates have been fetched and staged
    'completed',
    -- GBIF candidates could not be fetched due to an error
    'failed',
    -- GBIF candidates were fetched but none were found
    'no_candidates'
);

CREATE TABLE IF NOT EXISTS taxon_resolution (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    input_name CITEXT NOT NULL,
    input_authorship CITEXT,
    input_rank CITEXT,
    scientific_name CITEXT NOT NULL GENERATED ALWAYS AS (
        input_name || COALESCE(' ' || input_authorship, '')
    ) STORED,
    status resolution_status NOT NULL DEFAULT 'pending',
    gbif_status taxon_gbif_status NOT NULL DEFAULT 'pending',
    from_resolution_id UUID,
    sampling_target BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT taxon_resolution_unique UNIQUE (import_id, input_name),
    CONSTRAINT taxon_resolution_unique_id UNIQUE (import_id, id),
    CONSTRAINT taxon_resolution_from_resolution_fk FOREIGN KEY (import_id, from_resolution_id) REFERENCES taxon_resolution (import_id, id) ON DELETE CASCADE
);

CREATE TYPE taxon_match_type AS ENUM ('exact', 'fuzzy', 'name_only');
CREATE TYPE taxon_match_source AS ENUM ('internal', 'gbif', 'manual');


CREATE TABLE taxa_staging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    name CITEXT NOT NULL,
    authorship TEXT,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    parent_resolution_id UUID NOT NULL,
    FOREIGN KEY (import_id, parent_resolution_id) REFERENCES taxon_resolution(import_id, id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS taxon_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resolution_id UUID NOT NULL REFERENCES taxon_resolution (id) ON DELETE CASCADE,
    source taxon_match_source NOT NULL,
    match_type taxon_match_type NOT NULL,
    taxon_id UUID REFERENCES taxa(id) ON DELETE CASCADE,
    gbif_id INTEGER REFERENCES gbif_staging(key) ON DELETE CASCADE,
    staging_id UUID REFERENCES taxa_staging(id) ON DELETE CASCADE,
    score DOUBLE PRECISION,
    priority INTEGER NOT NULL,
    name CITEXT NOT NULL,
    authorship CITEXT,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    CONSTRAINT taxon_candidates_unique UNIQUE (resolution_id, id),
    CONSTRAINT taxon_candidates_target_check CHECK (
        (
            source = 'internal'
            AND taxon_id IS NOT NULL
        )
        OR (
            source = 'gbif'
            AND gbif_id IS NOT NULL
        )
        OR (
            source = 'manual'
            AND staging_id IS NOT NULL
        )
    )
);

-- TAXON RESOLUTION BACKLINK
ALTER TABLE taxon_resolution
ADD COLUMN IF NOT EXISTS resolved_candidate_id UUID,
    ADD CONSTRAINT taxon_resolution_resolved_candidate_fk FOREIGN KEY (id, resolved_candidate_id) REFERENCES taxon_candidates (resolution_id, id) ON DELETE
SET NULL (resolved_candidate_id),
    ADD CONSTRAINT taxon_resolution_resolved_status_check CHECK (
        status NOT IN ('auto_resolved', 'user_resolved')
        OR resolved_candidate_id IS NOT NULL
    );

-- A deleted or cleared selection must become eligible for resolution again.
CREATE FUNCTION taxon_resolution_reset_status() RETURNS trigger AS $$
BEGIN
    IF NEW.resolved_candidate_id IS NULL
        AND NEW.status IN ('auto_resolved', 'user_resolved') THEN
        NEW.status := 'pending';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER taxon_resolution_reset_status_trigger BEFORE
UPDATE OF resolved_candidate_id ON taxon_resolution
FOR EACH ROW EXECUTE FUNCTION taxon_resolution_reset_status();

CREATE UNIQUE INDEX taxon_candidates_internal_unique ON taxon_candidates (resolution_id, taxon_id)
WHERE source = 'internal';

CREATE UNIQUE INDEX taxon_candidates_gbif_unique ON taxon_candidates (resolution_id, gbif_id);
