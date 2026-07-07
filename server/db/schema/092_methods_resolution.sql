CREATE TYPE vocab_resolution_status AS ENUM (
    'auto',
    'selected',
    'pending',
    'request_creation',
    'discard'
);


CREATE TABLE IF NOT EXISTS sampling_methods_resolution (
    import_id UUID NOT NULL REFERENCES import_workflows (import_id) ON DELETE CASCADE,
    input_text TEXT NOT NULL,
    resolved_method_id UUID REFERENCES sampling_methods (id) ON DELETE CASCADE,
    status vocab_resolution_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (import_id, input_text)
)