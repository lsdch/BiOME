-- DATASETS
CREATE TABLE datasets (
	id ULID PRIMARY KEY,
	label TEXT NOT NULL,
	slug TEXT NOT NULL,
	description TEXT,
	pinned BOOLEAN NOT NULL DEFAULT FALSE,
	owner_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
	is_public BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT dataset_label_unique UNIQUE (label),
	CONSTRAINT dataset_slug_unique UNIQUE (slug),
	CONSTRAINT dataset_label_length CHECK (
		CHAR_LENGTH(BTRIM(label)) BETWEEN 4 AND 40
	)
);

-- DATASETS <-> OCCURRENCES
CREATE TABLE occurrences_datasets (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, occurrence_id)
);

CREATE INDEX occurrences_datasets_occurrence_idx ON occurrences_datasets (occurrence_id);

-- DATASETS <-> PUBLICATIONS
CREATE TABLE datasets_publications (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	publication_id UUID NOT NULL REFERENCES publications (id) ON DELETE RESTRICT,
	PRIMARY KEY (dataset_id, publication_id)
);

CREATE INDEX datasets_publications_publication_idx ON datasets_publications (publication_id);

-- DATASETS <-> USERS (curators)
CREATE TABLE datasets_curators (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, user_id)
);

CREATE INDEX datasets_curators_user_idx ON datasets_curators (user_id);

-- DATASETS <-> IMPORT_BATCHES
CREATE TABLE IF NOT EXISTS datasets_import_batches (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	import_batch_id UUID NOT NULL REFERENCES import_batches (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, import_batch_id)
);

CREATE INDEX datasets_import_batches_import_batch_idx ON datasets_import_batches (import_batch_id);