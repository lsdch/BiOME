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