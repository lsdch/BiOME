-- name: LockImportBatchForHabitats :one
SELECT id, status FROM import_batches WHERE id = @import_id FOR UPDATE;

-- name: InitSamplingHabitatsResolution :exec
WITH inputs AS (
    SELECT DISTINCT value::text AS input_text
    FROM import_samplings_occurrences i
        CROSS JOIN LATERAL unnest(i.habitats) AS value
    WHERE i.import_id = @import_id AND btrim(value::text) <> ''
)
INSERT INTO sampling_habitats_resolution (import_id, input_text, status)
SELECT @import_id, i.input_text,
    CASE WHEN h.id IS NULL THEN 'pending' ELSE 'auto_resolved' END::vocab_resolution_status
FROM inputs i LEFT JOIN habitats h ON h.label = i.input_text::citext
ON CONFLICT (import_id, input_text) DO NOTHING;

-- name: InitSamplingHabitatTargets :exec
INSERT INTO sampling_habitat_targets (import_id, input_text, habitat_id)
SELECT r.import_id, r.input_text, h.id
FROM sampling_habitats_resolution r
    JOIN habitats h ON h.label = r.input_text::citext
WHERE r.import_id = @import_id AND r.status = 'auto_resolved'
ON CONFLICT (import_id, input_text, habitat_id) DO NOTHING;

-- name: GetSamplingHabitatsResolution :many
SELECT * FROM sampling_habitats_resolution WHERE import_id = @import_id ORDER BY input_text;

-- name: GetSamplingHabitatTargets :many
SELECT * FROM sampling_habitat_targets WHERE import_id = @import_id ORDER BY input_text, id;

-- name: UpdateSamplingHabitatsResolution :execrows
UPDATE sampling_habitats_resolution SET status = @status
WHERE import_id = @import_id AND input_text = @input_text;

-- name: DeleteSamplingHabitatTargets :exec
DELETE FROM sampling_habitat_targets WHERE import_id = @import_id AND input_text = @input_text;

-- name: InsertSamplingHabitatTargets :exec
INSERT INTO sampling_habitat_targets (
    import_id, input_text, habitat_id, creation_name, creation_facet, creation_description
)
SELECT @import_id, @input_text, t.habitat_id, t.creation_name, t.creation_facet, t.creation_description
FROM jsonb_to_recordset(@targets::jsonb) AS t (
    habitat_id UUID, creation_name TEXT, creation_facet TEXT, creation_description TEXT
);

-- name: HasConflictingHabitatCreations :one
SELECT EXISTS (
    SELECT 1 FROM sampling_habitat_targets a
        JOIN sampling_habitat_targets b ON b.import_id = a.import_id
            AND b.creation_name = a.creation_name AND b.id <> a.id
    WHERE a.import_id = @import_id
        AND (a.creation_facet IS DISTINCT FROM b.creation_facet
            OR a.creation_description IS DISTINCT FROM b.creation_description)
)::boolean AS conflict;

-- name: CheckSamplingHabitatsReady :one
SELECT sampling_habitats_ready(@import_id)::boolean AS ready;

-- name: MaterializeRequestedHabitats :exec
INSERT INTO habitats (label, facet, description)
SELECT DISTINCT ON (t.creation_name) t.creation_name, t.creation_facet, t.creation_description
FROM sampling_habitat_targets t
    JOIN sampling_habitats_resolution r USING (import_id, input_text)
WHERE t.import_id = @import_id AND r.status = 'request_creation' AND t.creation_name IS NOT NULL
ORDER BY t.creation_name, t.id
ON CONFLICT (label) DO NOTHING;

-- name: ResolveCreatedHabitatTargets :exec
UPDATE sampling_habitat_targets t
SET habitat_id = h.id, creation_name = NULL, creation_facet = NULL, creation_description = NULL
FROM habitats h, sampling_habitats_resolution r
WHERE t.import_id = @import_id AND r.import_id = t.import_id AND r.input_text = t.input_text
    AND r.status = 'request_creation' AND h.label = t.creation_name;

-- name: MarkHabitatCreationsResolved :exec
UPDATE sampling_habitats_resolution r SET status = 'selected'
WHERE r.import_id = @import_id AND r.status = 'request_creation'
    AND NOT EXISTS (
        SELECT 1 FROM sampling_habitat_targets t
        WHERE t.import_id = r.import_id AND t.input_text = r.input_text AND t.habitat_id IS NULL
    );

-- name: MaterializeSamplingHabitats :exec
INSERT INTO samplings_habitats (sampling_id, habitat_id)
SELECT DISTINCT i.materialized_sampling_id, t.habitat_id
FROM import_samplings_occurrences i
    CROSS JOIN LATERAL unnest(i.habitats) AS value
    JOIN sampling_habitats_resolution r
        ON r.import_id = i.import_id AND r.input_text = value::text
    JOIN sampling_habitat_targets t ON t.import_id = r.import_id AND t.input_text = r.input_text
WHERE i.import_id = @import_id AND i.materialized_sampling_id IS NOT NULL
    AND r.status IN ('selected', 'auto_resolved') AND t.habitat_id IS NOT NULL
ON CONFLICT (sampling_id, habitat_id) DO NOTHING;

-- name: LockSamplingHabitatSelections :many
SELECT h.id FROM habitats h
WHERE EXISTS (
    SELECT 1 FROM sampling_habitat_targets t
    WHERE t.import_id = @import_id AND t.habitat_id = h.id
)
ORDER BY h.id FOR KEY SHARE OF h;

-- name: DeleteRedundantCreatedHabitatTargets :exec
DELETE FROM sampling_habitat_targets t
USING habitats h, sampling_habitats_resolution r
WHERE t.import_id = @import_id AND r.import_id = t.import_id AND r.input_text = t.input_text
    AND r.status = 'request_creation' AND t.creation_name = h.label
    AND EXISTS (
        SELECT 1 FROM sampling_habitat_targets other
        WHERE other.import_id = t.import_id AND other.input_text = t.input_text
            AND other.habitat_id = h.id
    );

-- name: HasUnmaterializedHabitatCreations :one
SELECT EXISTS (
    SELECT 1 FROM sampling_habitats_resolution
    WHERE import_id = @import_id AND status = 'request_creation'
)::boolean AS pending;

-- name: ReadImportBatchForHabitats :one
SELECT id FROM import_batches WHERE id = @import_id FOR SHARE;
