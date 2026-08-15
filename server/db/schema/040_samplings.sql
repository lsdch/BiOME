CREATE TABLE samplings (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	source_sampling_hash TEXT,
	comments CITEXT,
	-- LOCATION FIELDS
	site_code CITEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
	site_name CITEXT,
	site_locality CITEXT,
	site_country_code CHAR(3) REFERENCES countries (code),
	coordinates_precision INTEGER,
	coordinates geometry (Point, 4326) NOT NULL,
	latitude DOUBLE PRECISION NOT NULL GENERATED ALWAYS AS (ST_Y (coordinates)) STORED,
	longitude DOUBLE PRECISION NOT NULL GENERATED ALWAYS AS (ST_X (coordinates)) STORED,
	altitude INTEGER,
	-- DATE FIELDS
	event_date DATE,
	event_date_precision event_date_precision,
	-- METADATA FIELDS
	performed_by CITEXT [],
	duration INTEGER,
	access_points CITEXT [],
	import_batch_id UUID REFERENCES import_batches (id) ON DELETE
	SET NULL,
		-- R ~ 500m
		h3_index BIGINT NOT NULL GENERATED ALWAYS AS (
			h3_lat_lng_to_cell (
				point(ST_X (coordinates), ST_Y (coordinates)),
				12
			)::BIGINT
		) STORED,
		-- UTILITY FIELDS
		search_vector tsvector,
		-- CONSTRAINTS
		CONSTRAINT samplings_coordinates_precision_range CHECK (
			coordinates_precision IS NULL
			OR (
				coordinates_precision >= 0
				AND coordinates_precision <= 100000
			)
		),
		CONSTRAINT unique_sampling_hash_per_batch UNIQUE (source_sampling_hash, import_batch_id)
);


CREATE INDEX samplings_coordinates_gist_idx ON samplings USING GIST (coordinates);

CREATE INDEX samplings_h3_idx ON samplings (h3_index);

CREATE INDEX samplings_site_name_idx ON samplings (site_name text_pattern_ops);

-- TRIGGERS
CREATE FUNCTION samplings_search_vector_update () RETURNS trigger AS $$ --
BEGIN NEW.search_vector := to_tsvector('simple', coalesce(NEW.site_name, ''));
RETURN NEW;
END $$ LANGUAGE plpgsql;

CREATE TRIGGER samplings_search_vector_trigger BEFORE
INSERT
	OR
UPDATE ON samplings FOR EACH ROW EXECUTE FUNCTION samplings_search_vector_update ();

-- Association table linking samplings to taxa targets
CREATE TABLE sampling_target_taxa (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	taxon_id UUID NOT NULL REFERENCES taxa (id) ON DELETE RESTRICT,
	PRIMARY KEY (sampling_id, taxon_id)
);

CREATE INDEX sampling_target_taxa_taxon_idx ON sampling_target_taxa (taxon_id);
CREATE INDEX sampling_search_vector_idx ON samplings USING gin (search_vector);


CREATE VIEW samplings_with_country AS (
	SELECT s.*,
		c.code AS country_code,
		c.name AS country_name,
		c.continent AS country_continent,
		c.subcontinent AS country_subcontinent
	FROM samplings s
		LEFT JOIN countries c ON c.code = s.site_country_code
);