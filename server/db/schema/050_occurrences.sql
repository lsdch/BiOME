CREATE TABLE occurrences (
	-- ULID
	id CHAR(26) PRIMARY KEY,
	code TEXT NOT NULL,
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	type_status occurrence_type_status,
	comments TEXT,
	-- Identification fields
	taxon_id UUID NOT NULL REFERENCES taxa (id) ON DELETE RESTRICT,
	verbatim_identification TEXT,
	identified_by TEXT [],
	identification_date DATE,
	identification_date_precision event_date_precision,
	-- whether the identification is a confer (i.e. tentative) identification
	identification_confer BOOLEAN NOT NULL DEFAULT FALSE,
	identification_addendum TEXT,
	-- Content fields
	content_description TEXT,
	quantity_exact INTEGER,
	quantity_lower INTEGER,
	quantity_upper INTEGER,
	sources TEXT [],
	CONSTRAINT occurrence_quantity_shape CHECK (
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

CREATE INDEX occurrences_sampling_id_idx ON occurrences (sampling_id);

CREATE INDEX occurrences_type_status_idx ON occurrences (type_status);

CREATE TABLE IF NOT EXISTS occurrence_code_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	occurrence_id CHAR(26) NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	code TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);