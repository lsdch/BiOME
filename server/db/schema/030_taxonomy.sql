CREATE TABLE taxa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    gbif_id INTEGER UNIQUE,
    name CITEXT NOT NULL CONSTRAINT taxon_name_length CHECK (char_length(name) >= 2),
    scientific_name CITEXT NOT NULL GENERATED ALWAYS AS (TRIM(name || COALESCE(' ' || authorship, ''))) STORED,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    authorship CITEXT,
    -- RELATIONSHIPS
    accepted_taxon_id UUID REFERENCES taxa (id) ON DELETE
    SET NULL,
        parent_id UUID REFERENCES taxa (id) ON DELETE CASCADE,
        -- UTILITY FIELDS
        search_vector tsvector,
        -- CONSTRAINTS
        CONSTRAINT taxon_parent_required_for_non_kingdom CHECK (
            rank = 'kingdom'::taxon_rank
            OR parent_id IS NOT NULL
        ),
        CONSTRAINT taxon_synonym_requires_accepted_taxon CHECK (
            (
                status = 'synonym'::taxon_status
                AND accepted_taxon_id IS NOT NULL
            )
            OR (
                status <> 'synonym'::taxon_status
                AND accepted_taxon_id IS NULL
            )
        ),
        comments TEXT
);
CREATE INDEX taxa_name_rank_status_idx ON taxa (name, rank, status);
CREATE UNIQUE INDEX taxa_name_authorship_uidx ON taxa (name, COALESCE(authorship, ''));
CREATE INDEX taxa_parent_id_idx ON taxa (parent_id);
CREATE INDEX taxa_accepted_taxon_id_idx ON taxa (accepted_taxon_id);
-- Fuzzy search indexes
CREATE INDEX taxa_name_trgm_idx ON taxa USING gin (name gin_trgm_ops);
CREATE INDEX taxa_scientific_trgm_idx ON taxa USING gin (scientific_name gin_trgm_ops);
CREATE INDEX taxa_name_idx ON taxa (name text_pattern_ops);
CREATE INDEX taxa_search_vector_idx ON taxa USING gin (search_vector);

-- TRIGGERS
CREATE FUNCTION taxa_search_vector_update() RETURNS trigger AS $$ BEGIN NEW.search_vector := setweight(
    to_tsvector('simple', coalesce(NEW.scientific_name, '')),
    'A'
);
RETURN NEW;
END $$ LANGUAGE plpgsql;

CREATE TRIGGER taxa_search_vector_trigger BEFORE
INSERT
    OR
UPDATE ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_search_vector_update();


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