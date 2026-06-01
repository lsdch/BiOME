-- Minimal dataset and associations
CREATE TABLE datasets (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	import_hash TEXT NOT NULL UNIQUE,
	label TEXT NOT NULL,
	slug TEXT NOT NULL,
	description TEXT,
	submitted_by TEXT,
	assembled_by TEXT [],
	pinned BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT dataset_slug_unique UNIQUE (slug),
	CONSTRAINT dataset_label_length CHECK (
		CHAR_LENGTH(BTRIM(label)) BETWEEN 4 AND 40
	)
);

CREATE INDEX dataset_import_hash_idx ON datasets (import_hash);
CREATE INDEX dataset_slug_idx ON datasets (slug);

-- association dataset <-> occurrence
CREATE TABLE occurrences_datasets (
	dataset_id UUID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	occurrence_id CHAR(26) NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, occurrence_id)
);

CREATE INDEX occurrences_datasets_occurrence_idx ON occurrences_datasets (occurrence_id);

-- association dataset <-> article (publications)
CREATE TABLE datasets_publications (
	dataset_id UUID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	article_id UUID NOT NULL REFERENCES articles (id) ON DELETE RESTRICT,
	PRIMARY KEY (dataset_id, article_id)
);

CREATE INDEX datasets_publications_article_idx ON datasets_publications (article_id);

-- association dataset <-> users (curators)
CREATE TABLE dataset_curator (
	dataset_id UUID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, user_id)
);

CREATE INDEX dataset_curator_user_idx ON dataset_curator (user_id);