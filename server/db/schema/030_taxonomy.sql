CREATE TABLE taxa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    gbif_id INTEGER UNIQUE,
    name CITEXT NOT NULL CONSTRAINT taxon_name_length CHECK (char_length(name) >= 2),
    scientific_name CITEXT GENERATED ALWAYS AS (name || COALESCE(' ' || authorship, '')) STORED,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    authorship CITEXT,
    accepted_taxon_id UUID REFERENCES taxa (id) ON DELETE
    SET NULL,
        parent_id UUID REFERENCES taxa (id) ON DELETE CASCADE,
        CONSTRAINT taxon_parent_required_for_non_kingdom CHECK (
            rank = 'Kingdom'
            OR parent_id IS NOT NULL
        ),
        CONSTRAINT taxon_synonym_requires_accepted_taxon CHECK (
            (
                status = 'Synonym'
                AND accepted_taxon_id IS NOT NULL
            )
            OR (
                status <> 'Synonym'
                AND accepted_taxon_id IS NULL
            )
        ),
        comments TEXT
);
CREATE INDEX taxa_name_rank_status_idx ON taxa (name, rank, status);
CREATE UNIQUE INDEX taxa_name_authorship_uidx ON taxa (name, COALESCE(authorship, ''));
CREATE INDEX taxa_parent_id_idx ON taxa (parent_id);
CREATE INDEX taxa_accepted_taxon_id_idx ON taxa (accepted_taxon_id);
CREATE INDEX taxa_name_trgm_idx ON taxa USING gin (name gin_trgm_ops);
CREATE INDEX taxa_scientific_trgm_idx ON taxa USING gin (scientific_name gin_trgm_ops);