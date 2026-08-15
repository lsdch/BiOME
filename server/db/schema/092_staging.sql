CREATE TABLE import_samplings_occurrences (
    -- =========================
    -- INGESTION CONTEXT
    -- =========================
    id ULID PRIMARY KEY,
    import_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    row_number INTEGER NOT NULL,
    -- =========================
    -- SAMPLING
    -- =========================
    sampling_hash TEXT NOT NULL,
    sampling_comments TEXT,
    site_code TEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
    site_name TEXT,
    site_locality TEXT,
    site_country_code CHAR(3) REFERENCES countries (code),
    coordinates_precision INTEGER,
    longitude DOUBLE PRECISION NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    coordinates geometry (Point, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    ) STORED,
    altitude INTEGER,
    event_date DATE,
    event_date_precision event_date_precision,
    performed_by TEXT [],
    duration INTEGER,
    access_points CITEXT [],
    sampling_targets CITEXT [],
    sampling_fixatives CITEXT [],
    sampling_methods CITEXT [],
    habitats CITEXT [],
    -- =========================
    -- OCCURRENCE IDENTIFICATION
    -- =========================
    -- occurrence_code can be provided by the user, but will not be used as the actual occurrence code.
    -- Instead, the occurrence_code will be generated.
    -- The provided occurrence_code will be stored as a reference for the user, in the codes history.
    occurrence_code TEXT,
    generated_code TEXT,
    type_status occurrence_type_status,
    taxon_name CITEXT NOT NULL,
    taxon_authorship CITEXT,
    taxon_scientific_name CITEXT NOT NULL GENERATED ALWAYS AS (
        taxon_name || COALESCE(' ' || taxon_authorship, '')
    ) STORED,
    taxon_rank taxon_rank,
    verbatim_identification TEXT,
    identified_by TEXT [],
    identification_date DATE,
    identification_date_precision event_date_precision,
    identification_confer BOOLEAN NOT NULL DEFAULT FALSE,
    identification_addendum TEXT,
    -- =========================
    -- OCCURRENCE METADATA
    -- =========================
    content_description TEXT,
    quantity_exact INTEGER,
    quantity_lower INTEGER,
    quantity_upper INTEGER,
    sources TEXT [],
    occurrence_comments TEXT,
    -- =========================
    -- RESOLUTIONS
    -- =========================
    taxon_resolution_id UUID REFERENCES taxon_resolution (id) ON DELETE
    SET NULL,
        -- =========================
        -- MATERIALIZATIONS
        -- =========================
        materialized_sampling_id UUID REFERENCES samplings (id) ON DELETE
    SET NULL,
        materialized_occurrence_id ULID REFERENCES occurrences (id) ON DELETE
    SET NULL,
        -- =========================
        -- CONSTRAINTS
        -- =========================
        CONSTRAINT import_coordinate_precision_check CHECK (
            coordinates_precision IS NULL
            OR (
                coordinates_precision >= 0
                AND coordinates_precision <= 100000
            )
        ),
        CONSTRAINT import_sampling_occurrence_qty_check CHECK (
            (
                quantity_exact IS NULL
                AND quantity_lower IS NULL
                AND quantity_upper IS NULL
            )
            OR (
                quantity_exact IS NOT NULL
                AND quantity_lower IS NULL
                AND quantity_upper IS NULL
            )
            OR (
                quantity_exact IS NULL
                AND (
                    quantity_lower IS NOT NULL
                    OR quantity_upper IS NOT NULL
                )
                AND (
                    quantity_lower IS NULL
                    OR quantity_upper IS NULL
                    OR quantity_lower <= quantity_upper
                )
            )
        ) -- CONSTRAINT import_sampling_occurrence_collections_cardinality CHECK (
        --     collections IS NULL
        --     OR collection_vouchers IS NULL
        --     OR cardinality(collections) = cardinality(collection_vouchers)
        -- )
);

CREATE INDEX idx_staging_import ON import_samplings_occurrences(import_id);

CREATE INDEX idx_staging_hash ON import_samplings_occurrences(import_id, sampling_hash);

CREATE OR REPLACE VIEW samplings_staging AS
SELECT DISTINCT ON (import_id, sampling_hash) -- keep only the first row for each sampling hash within each import
    import_id,
    sampling_hash,
    row_number AS representative_row_number,
    site_code,
    coordinates,
    latitude,
    longitude,
    event_date,
    event_date_precision,
    site_name,
    site_locality,
    site_country_code,
    coordinates_precision,
    altitude,
    performed_by,
    duration,
    access_points,
    sampling_targets,
    sampling_fixatives,
    sampling_methods,
    habitats,
    imported_at,
    sampling_comments,
    materialized_sampling_id
FROM import_samplings_occurrences
ORDER BY import_id,
    sampling_hash,
    row_number;


CREATE TABLE IF NOT EXISTS collections_staging (
    occurrence_id ULID NOT NULL REFERENCES import_samplings_occurrences (id) ON DELETE CASCADE,
    collection_name TEXT NOT NULL,
    vouchers TEXT [],
    PRIMARY KEY (occurrence_id, collection_name)
);