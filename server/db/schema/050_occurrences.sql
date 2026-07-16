CREATE TABLE IF NOT EXISTS occurrences (
	-- ULID
	id ULID PRIMARY KEY,
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
	-- Metadata fields
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
	import_batch_id ULID REFERENCES import_batches (id) ON DELETE
	SET NULL,
		-- Constraints
		CONSTRAINT occurrence_quantity_shape CHECK (
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
		)
);

CREATE INDEX occurrences_sampling_id_idx ON occurrences (sampling_id);

CREATE INDEX occurrences_type_status_idx ON occurrences (type_status);

CREATE INDEX occurrences_code_idx ON occurrences (code text_pattern_ops);

CREATE TRIGGER occurrences_set_updated_at BEFORE
UPDATE ON occurrences FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS occurrence_code_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	code TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE VIEW occurrence_codes_to_update AS with codes_projection AS (
	SELECT o.id,
		o.import_batch_id,
		o.code AS current_code,
		format(
			'%s[%s|%s]',
			REPLACE(t.name, ' ', '_'),
			COALESCE(
				s.site_code,
				format('%sN,%sE', s.latitude, s.longitude)
			),
			(
				CASE
					WHEN s.event_date IS NOT NULL THEN (
						CASE
							WHEN s.event_date_precision = 'day'
							OR s.event_date_precision IS NULL THEN to_char(s.event_date, 'YYYY-MM-DD')
							WHEN s.event_date_precision = 'month' THEN to_char(s.event_date, 'YYYY-MM')
							WHEN s.event_date_precision = 'year' THEN to_char(s.event_date, 'YYYY')
						END
					)
					ELSE 'N/A'
				END
			)
		) AS computed_code
	FROM occurrences o
		JOIN samplings s ON s.id = o.sampling_id
		JOIN taxa t ON t.id = o.taxon_id
)
SELECT *
FROM codes_projection
WHERE current_code IS DISTINCT
FROM computed_code;