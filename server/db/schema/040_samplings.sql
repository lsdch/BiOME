CREATE TABLE samplings (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	sampling_hash TEXT NOT NULL UNIQUE,
	notes TEXT,
	-- LOCATION FIELDS
	site_code TEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
	site_name TEXT,
	site_locality TEXT,
	site_country_code CHAR(3) REFERENCES countries (code),
	coordinates_precision INTEGER,
	coordinates geometry (Point, 4326) NOT NULL,
	latitude REAL NOT NULL GENERATED ALWAYS AS (ST_Y (coordinates)) STORED,
	longitude REAL NOT NULL GENERATED ALWAYS AS (ST_X (coordinates)) STORED,
	altitude INTEGER,
	-- DATE FIELDS
	event_date DATE,
	event_date_precision event_date_precision,
	-- METADATA FIELDS
	performed_by TEXT [],
	duration INTEGER,
	access_points TEXT [],
	-- R ~ 500m
	h3_res8 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 8)
	) STORED,
	-- R ~ 1km
	h3_res7 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 7)
	) STORED,
	-- R ~ 3 km
	h3_res6 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 6)
	) STORED,
	-- R ~ 10km
	h3_res5 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 5)
	) STORED,
	-- R ~ 25km
	h3_res4 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 4)
	) STORED,
	-- R ~ 60km
	h3_res3 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 3)
	) STORED,
	-- R ~ 166km
	h3_res2 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 2)
	) STORED,
	CONSTRAINT samplings_coordinates_precision_range CHECK (
		coordinates_precision IS NULL
		OR (
			coordinates_precision >= 0
			AND coordinates_precision <= 100000
		)
	)
);

CREATE INDEX idx_samplings_hash ON samplings(sampling_hash);
CREATE INDEX samplings_coordinates_gist_idx ON samplings USING GIST (coordinates);
CREATE INDEX samplings_h3_r8_idx ON samplings (h3_res8);
CREATE INDEX samplings_h3_r7_idx ON samplings (h3_res7);
CREATE INDEX samplings_h3_r6_idx ON samplings (h3_res6);
CREATE INDEX samplings_h3_r5_idx ON samplings (h3_res5);
CREATE INDEX samplings_h3_r4_idx ON samplings (h3_res4);
CREATE INDEX samplings_h3_r3_idx ON samplings (h3_res3);
CREATE INDEX samplings_h3_r2_idx ON samplings (h3_res2);

-- Association table linking samplings to taxa targets
CREATE TABLE sampling_target_taxa (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	taxon_id UUID NOT NULL REFERENCES taxa (id) ON DELETE RESTRICT,
	PRIMARY KEY (sampling_id, taxon_id)
);

CREATE INDEX sampling_target_taxa_taxon_idx ON sampling_target_taxa (taxon_id);