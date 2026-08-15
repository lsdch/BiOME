CREATE TYPE import_batch_status AS ENUM (
    'created',
    'staged',
    'completed',
    'failed',
    'canceled'
);

CREATE TABLE IF NOT EXISTS import_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label TEXT NOT NULL,
    description TEXT,
    status import_batch_status NOT NULL DEFAULT 'created',
    assembled_by TEXT [],
    created_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    completed_by UUID REFERENCES users (id) ON DELETE RESTRICT,
    taxonomic_scope INT NOT NULL REFERENCES gbif_staging (key) ON DELETE RESTRICT,
    CONSTRAINT import_batches_label_length CHECK (
        CHAR_LENGTH(BTRIM(label)) BETWEEN 3 AND 40
    )
);