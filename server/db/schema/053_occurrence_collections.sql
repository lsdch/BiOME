CREATE TABLE occurrence_collections (
	collection_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	vouchers TEXT [],
	CONSTRAINT occurrence_collections_name_not_empty CHECK (char_length(btrim(name)) >= 2)
);
CREATE INDEX occurrence_collections_name_idx ON occurrence_collections (name);