CREATE TABLE import_workflows (
    import_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    label TEXT NOT NULL,
    description TEXT,
    assembled_by TEXT [],
    created_by UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);

ALTER TABLE import_batches
ADD COLUMN IF NOT EXISTS workflow_id UUID REFERENCES import_workflows (import_id) ON DELETE
SET NULL;