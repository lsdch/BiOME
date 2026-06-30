-- Minimal dataset and associations
CREATE TABLE datasets (
	id ULID PRIMARY KEY DEFAULT gen_random_uuid (),
	label TEXT NOT NULL,
	slug TEXT NOT NULL,
	description TEXT,
	pinned BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT dataset_slug_unique UNIQUE (slug),
	CONSTRAINT dataset_label_length CHECK (
		CHAR_LENGTH(BTRIM(label)) BETWEEN 4 AND 40
	)
);

CREATE INDEX dataset_slug_idx ON datasets (slug);

-- association dataset <-> occurrence
CREATE TABLE occurrences_datasets (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, occurrence_id)
);

CREATE INDEX occurrences_datasets_occurrence_idx ON occurrences_datasets (occurrence_id);

-- association dataset <-> article (publications)
CREATE TABLE datasets_publications (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	article_id UUID NOT NULL REFERENCES articles (id) ON DELETE RESTRICT,
	PRIMARY KEY (dataset_id, article_id)
);

CREATE INDEX datasets_publications_article_idx ON datasets_publications (article_id);

-- association dataset <-> users (curators)
CREATE TABLE datasets_curators (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, user_id)
);

CREATE INDEX datasets_curators_user_idx ON datasets_curators (user_id);