CREATE TABLE countries (
	code TEXT PRIMARY KEY CHECK (code ~ '^[A-Z]{3}$'),
	name TEXT NOT NULL UNIQUE,
	continent TEXT NOT NULL,
	subcontinent TEXT NOT NULL
);