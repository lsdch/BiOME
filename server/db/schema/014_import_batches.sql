CREATE TABLE IF NOT EXISTS import_batches (
    id ULID PRIMARY KEY,
    label TEXT NOT NULL,
    description TEXT,
    assembled_by TEXT [],
    created_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Added after declaration of import_workflows table
    -- workflow_id UUID REFERENCES import_workflows (import_id) ON DELETE
    -- SET NULL,
    CONSTRAINT import_batches_label_length CHECK (
        CHAR_LENGTH(BTRIM(label)) BETWEEN 4 AND 40
    )
);