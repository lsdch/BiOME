CREATE TABLE IF NOT EXISTS gbif_staging (
    key INTEGER PRIMARY KEY,
    parent TEXT,
    parent_key INTEGER,
    canonical_name TEXT NOT NULL,
    scientific_name TEXT NOT NULL,
    status TEXT NOT NULL,
    rank TEXT NOT NULL,
    name_type TEXT NOT NULL,
    kingdom_key INTEGER,
    phylum_key INTEGER,
    class_key INTEGER,
    order_key INTEGER,
    family_key INTEGER,
    genus_key INTEGER,
    species_key INTEGER,
    higher_taxon_keys INTEGER [],
    higher_taxon_names TEXT [],
    authorship text,
    num_descendants INTEGER,
    accepted_key INTEGER,
    accepted_name TEXT
);

CREATE TABLE IF NOT EXISTS gbif_dependencies (
    import_hash TEXT NOT NULL,
    key INTEGER NOT NULL,
    PRIMARY KEY (import_hash, key)
);

CREATE TYPE taxon_match_type AS ENUM ('exact', 'fuzzy', 'name_only');
CREATE TYPE taxon_match_source AS ENUM ('internal', 'gbif', 'manual');

CREATE TABLE IF NOT EXISTS taxon_candidates (
    import_hash TEXT NOT NULL,
    input_name CITEXT NOT NULL,
    source taxon_match_source NOT NULL,
    match_type taxon_match_type NOT NULL,
    taxon_id UUID REFERENCES taxa(id),
    gbif_id INTEGER REFERENCES gbif_staging(key),
    score DOUBLE PRECISION,
    priority INTEGER NOT NULL,
    name CITEXT NOT NULL,
    authorship CITEXT,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    CONSTRAINT taxon_candidates_target_check CHECK (
        (
            source = 'internal'
            AND taxon_id IS NOT NULL
        )
        OR (
            source = 'gbif'
            AND gbif_id IS NOT NULL
            AND taxon_id IS NULL
        )
    )
);

CREATE UNIQUE INDEX taxon_candidates_internal_unique ON taxon_candidates (import_hash, input_name, taxon_id)
WHERE source = 'internal';

CREATE UNIQUE INDEX taxon_candidates_gbif_unique ON taxon_candidates (import_hash, input_name, gbif_id)
WHERE source = 'gbif';

CREATE TABLE taxa_staging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_hash TEXT NOT NULL,
    name CITEXT NOT NULL,
    authorship TEXT,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    parent_source taxon_match_source NOT NULL,
    parent_taxa_id UUID REFERENCES taxa(id),
    parent_gbif_id INTEGER REFERENCES gbif_staging(key),
    parent_input_name TEXT,
    CHECK (
        (
            parent_source = 'internal'
            AND parent_taxa_id IS NOT NULL
        )
        OR (
            parent_source = 'gbif'
            AND parent_gbif_id IS NOT NULL
        )
        OR (
            parent_source = 'manual'
            AND parent_input_name IS NOT NULL
        )
    )
);

CREATE TYPE resolution_status AS ENUM (
    'pending',
    'auto_resolved',
    'user_resolved',
    'needs_decision'
);

CREATE TYPE taxon_gbif_status AS ENUM (
    'skipped',
    'pending',
    'completed',
    'failed'
);

CREATE TABLE IF NOT EXISTS taxon_resolution (
    import_hash TEXT NOT NULL,
    input_name CITEXT NOT NULL,
    source taxon_match_source,
    gbif_id INTEGER REFERENCES gbif_staging(key),
    taxon_id UUID REFERENCES taxa(id),
    staging_id UUID REFERENCES taxa_staging(id),
    status resolution_status DEFAULT 'pending',
    gbif_status taxon_gbif_status DEFAULT 'skipped',
    PRIMARY KEY (import_hash, input_name)
);