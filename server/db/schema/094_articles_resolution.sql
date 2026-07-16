CREATE TABLE publications_staging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- =========================
    -- INGESTION CONTEXT
    -- =========================
    import_id UUID NOT NULL REFERENCES import_workflows (import_id) ON DELETE CASCADE,
    occurrence_row_number INTEGER NOT NULL,
    -- =========================
    -- ARTICLE METADATA
    -- =========================
    doi TEXT,
    authors TEXT [],
    year INTEGER,
    title TEXT,
    journal TEXT,
    verbatim TEXT
);


CREATE TYPE publication_resolution_type AS ENUM (
    'crossref',
    'doi',
    'verbatim',
    'manual'
);
CREATE TYPE publication_resolution_status AS ENUM (
    'pending',
    'resolved',
    'failed',
    'manual_required'
);

CREATE TABLE publication_resolution (
    staging_id UUID NOT NULL REFERENCES publications_staging (id) ON DELETE CASCADE,
    PRIMARY KEY (staging_id),
    publication_id UUID REFERENCES articles (id) ON DELETE
    SET NULL,
        resolution_type publication_resolution_type,
        status publication_resolution_status NOT NULL DEFAULT 'pending',
        crossref_payload JSONB
)