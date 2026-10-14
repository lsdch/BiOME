CREATE TYPE vocab_resolution_status AS ENUM (
    'auto_resolved',
    'selected',
    'pending',
    'request_creation',
    'discard'
);


CREATE TABLE IF NOT EXISTS sampling_methods_resolution (
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    input_text TEXT NOT NULL,
    resolved_method_id UUID REFERENCES sampling_methods (id) ON DELETE CASCADE,
    status vocab_resolution_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (import_id, input_text)
); 

CREATE TABLE IF NOT EXISTS sampling_fixatives_resolution (
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    input_text TEXT NOT NULL,
    resolved_fixative_id UUID REFERENCES fixatives (id) ON DELETE CASCADE,
    status vocab_resolution_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (import_id, input_text)
);

CREATE TABLE sampling_target_resolution (
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    sampling_hash TEXT NOT NULL,
    resolution_id UUID NOT NULL REFERENCES taxon_resolution (id) ON DELETE CASCADE,
    PRIMARY KEY (import_id, sampling_hash, resolution_id)
);

CREATE TABLE sampling_habitats_resolution (
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    input_text TEXT NOT NULL CHECK (btrim(input_text) <> ''),
    status vocab_resolution_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (import_id, input_text)
);

CREATE TABLE sampling_habitat_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_id UUID NOT NULL,
    input_text TEXT NOT NULL,
    habitat_id UUID REFERENCES habitats (id) ON DELETE SET NULL,
    creation_name CITEXT,
    creation_facet CITEXT,
    creation_description TEXT,
    FOREIGN KEY (import_id, input_text)
        REFERENCES sampling_habitats_resolution (import_id, input_text) ON DELETE CASCADE,
    UNIQUE (import_id, input_text, habitat_id),
    UNIQUE (import_id, input_text, creation_name),
    -- A NULL habitat with no creation is retained if its selected habitat was deleted.
    -- The trigger below makes that resolution pending again.
    CHECK (creation_name IS NULL OR (habitat_id IS NULL AND btrim(creation_name) <> '')),
    CHECK (creation_name IS NOT NULL OR (creation_facet IS NULL AND creation_description IS NULL)),
    CHECK (creation_facet IS NULL OR btrim(creation_facet) <> ''),
    CHECK (char_length(coalesce(creation_description, '')) <= 4000)
);
CREATE INDEX sampling_habitat_targets_habitat_idx ON sampling_habitat_targets (habitat_id);

CREATE FUNCTION reset_deleted_habitat_resolution() RETURNS trigger AS $$
BEGIN
    IF OLD.habitat_id IS NOT NULL AND NEW.habitat_id IS NULL THEN
        UPDATE sampling_habitats_resolution SET status = 'pending'
        WHERE import_id = NEW.import_id AND input_text = NEW.input_text;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER sampling_habitat_deleted
AFTER UPDATE OF habitat_id ON sampling_habitat_targets
FOR EACH ROW EXECUTE FUNCTION reset_deleted_habitat_resolution();

-- Check the source values too: a missing resolution must never silently lose habitats.
CREATE FUNCTION sampling_habitats_ready(batch_id UUID) RETURNS BOOLEAN AS $$
    WITH inputs AS (
        SELECT DISTINCT value::text AS input_text
        FROM import_samplings_occurrences i
            CROSS JOIN LATERAL unnest(i.habitats) AS value
        WHERE i.import_id = batch_id AND btrim(value::text) <> ''
    ), targets AS (
        SELECT input_text, count(*) AS total,
            count(*) FILTER (WHERE habitat_id IS NULL AND creation_name IS NULL) AS missing,
            count(*) FILTER (WHERE creation_name IS NOT NULL) AS creations
        FROM sampling_habitat_targets
        WHERE import_id = batch_id
        GROUP BY input_text
    )
    SELECT NOT EXISTS (
        SELECT 1 FROM inputs i
            LEFT JOIN sampling_habitats_resolution r
                ON r.import_id = batch_id AND r.input_text = i.input_text
            LEFT JOIN targets t ON t.input_text = i.input_text
        WHERE r.status IS NULL OR (
            (r.status = 'discard' AND coalesce(t.total, 0) = 0)
            OR (r.status IN ('auto_resolved', 'selected') AND t.total > 0
                AND t.missing = 0 AND t.creations = 0)
            OR (r.status = 'request_creation' AND t.total > 0
                AND t.missing = 0 AND t.creations > 0)
        ) IS NOT TRUE
    );
$$ LANGUAGE SQL STABLE;
