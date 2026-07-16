CREATE TABLE IF NOT EXISTS sampling_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    code CITEXT NOT NULL UNIQUE,
    name CITEXT NOT NULL UNIQUE,
    description TEXT
);


CREATE TABLE IF NOT EXISTS events_sampling_methods (
    sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
    method_id UUID NOT NULL REFERENCES sampling_methods (id) ON DELETE CASCADE,
    PRIMARY KEY (sampling_id, method_id)
);