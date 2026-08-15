-- Publications (bibliographic references)
CREATE TABLE publications (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	authors TEXT [],
	year INTEGER,
	title TEXT,
	journal TEXT,
	verbatim TEXT NOT NULL,
	doi DOI,
	comments TEXT,
	CONSTRAINT doi_unique UNIQUE (doi)
);

CREATE INDEX publications_year_idx ON publications (year);

-- Association table linking occurrences to publications (published_in)
CREATE TABLE occurrences_publications (
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	publication_id UUID NOT NULL REFERENCES publications (id) ON DELETE CASCADE,
	PRIMARY KEY (occurrence_id, publication_id)
);