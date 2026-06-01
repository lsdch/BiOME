-- Articles (bibliographic references)
CREATE TABLE articles (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	authors TEXT [] NOT NULL,
	year INTEGER NOT NULL,
	title TEXT,
	journal TEXT,
	verbatim TEXT,
	doi TEXT,
	comments TEXT,
	code TEXT NOT NULL,
	CONSTRAINT article_code_unique UNIQUE (code)
);

CREATE INDEX article_code_year_idx ON articles (code, year);

-- Association table linking occurrences to articles (published_in)
CREATE TABLE occurrence_article (
	occurrence_id CHAR(26) NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	article_id UUID NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
	PRIMARY KEY (occurrence_id, article_id)
);

CREATE INDEX occurrence_article_article_idx ON occurrence_article (article_id);