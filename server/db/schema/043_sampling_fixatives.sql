CREATE TABLE IF NOT EXISTS fixatives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    code CITEXT NOT NULL UNIQUE,
    name CITEXT NOT NULL UNIQUE,
    description TEXT
);


CREATE TABLE IF NOT EXISTS samplings_fixatives (
    sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
    fixative_id UUID NOT NULL REFERENCES fixatives (id) ON DELETE CASCADE,
    PRIMARY KEY (sampling_id, fixative_id)
);