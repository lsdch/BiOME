-- Users table
CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	login CITEXT NOT NULL UNIQUE,
	email CITEXT NOT NULL,
	password_hash TEXT NOT NULL,
	role user_role NOT NULL DEFAULT 'visitor',
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	organisation TEXT,
	contact TEXT,
	bio TEXT,
	full_name TEXT NOT NULL GENERATED ALWAYS AS (first_name || ' ' || last_name) STORED,
	active BOOLEAN NOT NULL DEFAULT true,
	email_verified_at TIMESTAMPTZ,
	CONSTRAINT users_login_unique UNIQUE (login),
	CONSTRAINT users_email_unique UNIQUE (email),
	CONSTRAINT users_login_length CHECK (char_length(btrim(login)) >= 2),
	CONSTRAINT users_first_name_length CHECK (char_length(btrim(first_name)) >= 2),
	CONSTRAINT users_last_name_length CHECK (char_length(btrim(last_name)) >= 2)
);

CREATE INDEX users_full_name_idx ON users (full_name);

-- Request-scoped helpers.
-- The application can set `biome.current_user_id` for the current transaction/session.
CREATE OR REPLACE FUNCTION current_request_user_id () RETURNS UUID LANGUAGE SQL STABLE AS $$
SELECT nullif(
		current_setting('biome.current_user_id', true),
		''
	)::uuid;
$$;

CREATE OR REPLACE FUNCTION current_request_user () RETURNS users LANGUAGE SQL STABLE AS $$
SELECT *
FROM users
WHERE id = current_request_user_id();
$$;