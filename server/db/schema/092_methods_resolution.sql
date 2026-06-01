CREATE TYPE method_resolution_status AS ENUM (
    'auto',
    'selected',
    'pending',
    'request_creation',
    'discard'
);


CREATE TABLE IF NOT EXISTS sampling_methods_resolution (
    import_hash TEXT NOT NULL,
    input_text TEXT NOT NULL,
    resolved_method_id UUID REFERENCES sampling_methods (id) ON DELETE CASCADE,
    status method_resolution_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (import_hash, input_text)
)