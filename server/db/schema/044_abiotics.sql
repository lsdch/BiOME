CREATE TABLE IF NOT EXISTS abiotic_params (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name CITEXT NOT NULL,
    code CITEXT NOT NULL UNIQUE,
    description TEXT,
    unit TEXT NOT NULL,
    CONSTRAINT abiotic_params_name_length CHECK (
        CHAR_LENGTH(BTRIM(name)) BETWEEN 4 AND 40
    )
);

CREATE TABLE IF NOT EXISTS abiotic_measurements (
    sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
    param_id UUID NOT NULL REFERENCES abiotic_params (id) ON DELETE RESTRICT,
    value NUMERIC NOT NULL,
    CONSTRAINT abiotic_measurements_pkey PRIMARY KEY (sampling_id, param_id)
)