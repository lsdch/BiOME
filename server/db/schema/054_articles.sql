-- Articles (bibliographic references)
CREATE TABLE articles (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	authors TEXT [] NOT NULL,
	year INTEGER NOT NULL,
	title TEXT,
	journal TEXT,
	verbatim TEXT,
	doi TEXT,
	comments TEXT
);

CREATE INDEX articles_year_idx ON articles (year);

-- Association table linking occurrences to articles (published_in)
CREATE TABLE occurrences_articles (
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	article_id UUID NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
	PRIMARY KEY (occurrence_id, article_id)
);