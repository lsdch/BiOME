CREATE TABLE IF NOT EXISTS import_batches (
    id ULID PRIMARY KEY,
    label TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    submitted_by TEXT,
    assembled_by TEXT [],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT import_batches_label_length CHECK (
        CHAR_LENGTH(BTRIM(label)) BETWEEN 4 AND 40
    ),
    CONSTRAINT import_batches_slug_length CHECK (
        CHAR_LENGTH(BTRIM(slug)) BETWEEN 4 AND 40
    )
);