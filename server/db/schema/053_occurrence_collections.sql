CREATE TABLE occurrence_collections (
	occurrence_id CHAR(26) NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	vouchers TEXT [],
	CONSTRAINT occurrence_collections_name_not_empty CHECK (char_length(btrim(name)) >= 2),
	PRIMARY KEY (occurrence_id, name)
);
CREATE INDEX occurrence_collections_name_idx ON occurrence_collections (name);