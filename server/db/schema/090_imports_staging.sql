CREATE TYPE gbif_import_status AS ENUM (
    'pending',
    'in_progress',
    'completed',
    'failed'
);

CREATE TABLE import_workflows (
    import_hash TEXT PRIMARY KEY,
    label TEXT,
    gbif_status gbif_import_status NOT NULL DEFAULT 'pending',
    gbif_candidates_total INTEGER,
    gbif_candidates_fetched INTEGER,
    gbif_claimed_at TIMESTAMPTZ,
    gbif_updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);


CREATE TABLE import_samplings_occurrences (
    -- =========================
    -- INGESTION CONTEXT
    -- =========================
    import_hash TEXT NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    row_number INTEGER NOT NULL,
    PRIMARY KEY (import_hash, row_number),
    -- =========================
    -- SAMPLING
    -- =========================
    sampling_hash TEXT NOT NULL,
    sampling_comments TEXT,
    site_code TEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
    site_name TEXT,
    site_locality TEXT,
    site_country_code CHAR(3) NOT NULL REFERENCES countries (code),
    coordinates_precision INTEGER,
    longitude REAL NOT NULL,
    latitude REAL NOT NULL,
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
    type_status occurrence_type_status,
    taxon_name CITEXT NOT NULL,
    taxon_authorship CITEXT,
    taxon_scientific_name CITEXT NOT NULL GENERATED ALWAYS AS (
        taxon_name || COALESCE(' ' || taxon_authorship, '')
    ) STORED,
    taxon_rank CITEXT,
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
    )
);

CREATE INDEX idx_staging_import ON import_samplings_occurrences(import_hash);

CREATE INDEX idx_staging_hash ON import_samplings_occurrences(import_hash, sampling_hash);