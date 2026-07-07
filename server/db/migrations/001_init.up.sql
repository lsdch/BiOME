CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS migrations;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE EXTENSION IF NOT EXISTS h3;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- used for case-insensitive text type and trigram indexes
CREATE EXTENSION IF NOT EXISTS citext;

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE TYPE coordinates_precision AS ENUM('<100m', '<1km', '<10km', '10-100km');

CREATE TYPE event_date_precision AS ENUM('day', 'month', 'year');

CREATE TYPE taxon_rank AS ENUM(
	'Subspecies',
	'Species',
	'Subgenus',
	'Genus',
	'Family',
	'Order',
	'Class',
	'Phylum',
	'Kingdom'
);

CREATE TYPE taxon_status AS ENUM(
	'Accepted',
	'Synonym',
	'Doubtful',
	'Unreferenced',
	'Unclassified'
);

CREATE TYPE occurrence_type_status AS ENUM('Holotype', 'Neotype', 'Topotype');

-- Roles for platform users (aligned with EdgeDB `people::UserRole`)
CREATE TYPE user_role AS ENUM('Visitor', 'Contributor', 'Maintainer', 'Admin');
CREATE TABLE countries (
	code TEXT PRIMARY KEY CHECK (code ~ '^[A-Z]{3}$'),
	name TEXT NOT NULL UNIQUE,
	continent TEXT NOT NULL,
	subcontinent TEXT NOT NULL
);
CREATE TYPE vocab_resolution AS (
	id uuid,
	code text,
	label text,
	confidence float,
	match_type text,
	should_accept boolean
);

CREATE OR REPLACE FUNCTION vocab_decision(
		v_id uuid,
		v_code text,
		v_label text,
		v_score float,
		min_auto float
	) RETURNS vocab_resolution LANGUAGE sql AS $$
SELECT v_id,
	v_code,
	v_label,
	COALESCE(v_score, 0),
	CASE
		WHEN v_score = 1 THEN 'exact'
		WHEN v_score >= min_auto THEN 'fuzzy'
		ELSE 'none'
	END,
	v_score >= min_auto;
$$;
-- Users table
CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	login CITEXT NOT NULL UNIQUE,
	email CITEXT NOT NULL,
	password_hash TEXT NOT NULL,
	role user_role NOT NULL DEFAULT 'Visitor',
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	organisation TEXT,
	contact TEXT,
	comments TEXT,
	full_name TEXT GENERATED ALWAYS AS (first_name || ' ' || last_name) STORED,
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
-- Invitations and tokens (separate metadata and token storage)
-- This design separates invitation metadata from token hashes and allows multiple tokens per invitation (rotation).
CREATE TYPE invitation_status AS ENUM('pending', 'redeemed', 'cancelled', 'expired');
CREATE TABLE invitations (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	email CITEXT NOT NULL,
	invitee_name TEXT,
	organisation TEXT,
	role user_role NOT NULL DEFAULT 'Visitor',
	message TEXT,
	inviter_id UUID REFERENCES users (id) ON DELETE
	SET NULL,
		status invitation_status NOT NULL DEFAULT 'pending',
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		expires_at TIMESTAMPTZ NOT NULL,
		redeemed_at TIMESTAMPTZ,
		revoked_at TIMESTAMPTZ,
		revoked_by UUID REFERENCES users (id) ON DELETE
	SET NULL,
		CONSTRAINT invitations_dates CHECK (expires_at > created_at)
);

-- Only one pending invitation per email
CREATE UNIQUE INDEX invitations_one_active_per_email_idx ON invitations (email)
WHERE status = 'pending';

CREATE INDEX invitations_expires_idx ON invitations (expires_at);
CREATE INDEX invitations_inviter_idx ON invitations (inviter_id);


-- Tokens table: store only a hash of the raw token for verification.
-- Redeeming an invitation implicitly verifies the invited email.
CREATE TABLE invitation_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	invitation_id UUID NOT NULL REFERENCES invitations (id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	consumed BOOLEAN NOT NULL DEFAULT false,
	consumed_by UUID REFERENCES users (id) ON DELETE
	SET NULL,
		consumed_at TIMESTAMPTZ,
		CONSTRAINT invitation_tokens_hash_unique UNIQUE (token_hash),
		CONSTRAINT invitation_tokens_consumed_time CHECK (
			(consumed = false)
			OR (consumed_at IS NOT NULL)
		)
);
CREATE INDEX invitation_tokens_hash_idx ON invitation_tokens (token_hash);
CREATE INDEX invitation_tokens_invitation_idx ON invitation_tokens (invitation_id);
-- Utility hashing function (uses pgcrypto.digest). Applications should pass only the raw token to the hashing function
CREATE OR REPLACE FUNCTION token_sha256 (txt TEXT) RETURNS TEXT LANGUAGE SQL IMMUTABLE AS $$
SELECT encode(digest(txt, 'sha256'), 'hex');
$$;
-- Example usage (application side):
-- INSERT INTO invitations (email, invitee_name, organisation, role, message, expires_at)
--   VALUES (..., now() + interval '7 days') RETURNING id;
-- INSERT INTO invitation_tokens (invitation_id, token_hash)
--   VALUES (:inv_id, token_sha256(:raw_token));-- User account requests and email-change requests.
-- Both flows require explicit email verification; invitation-based accounts bypass these tables.
CREATE TYPE user_account_request_status AS ENUM('pending', 'verified', 'cancelled', 'expired');
CREATE TYPE user_email_change_request_status AS ENUM(
	'pending',
	'verified',
	'applied',
	'cancelled',
	'expired'
);

CREATE TABLE user_account_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	email CITEXT NOT NULL,
	name TEXT NOT NULL,
	motivations TEXT,
	status user_account_request_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expires_at TIMESTAMPTZ NOT NULL,
	verified_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT user_account_requests_dates CHECK (expires_at > created_at),
	CONSTRAINT user_account_requests_verified_time CHECK (
		verified_at IS NULL
		OR verified_at >= created_at
	),
	CONSTRAINT user_account_requests_cancelled_time CHECK (
		cancelled_at IS NULL
		OR cancelled_at >= created_at
	)
);

CREATE UNIQUE INDEX user_account_requests_one_active_per_email_idx ON user_account_requests (email)
WHERE status = 'pending';

CREATE INDEX user_account_requests_expires_idx ON user_account_requests (expires_at);

-- Tokens table: store only a hash of the raw token for verification.
CREATE TABLE user_account_request_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_account_request_id UUID NOT NULL REFERENCES user_account_requests (id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	consumed BOOLEAN NOT NULL DEFAULT false,
	consumed_at TIMESTAMPTZ,
	CONSTRAINT user_account_request_tokens_hash_unique UNIQUE (token_hash),
	CONSTRAINT user_account_request_tokens_consumed_time CHECK (
		(consumed = false)
		OR (consumed_at IS NOT NULL)
	)
);

CREATE INDEX user_account_request_tokens_hash_idx ON user_account_request_tokens (token_hash);
CREATE INDEX user_account_request_tokens_request_idx ON user_account_request_tokens (user_account_request_id);


CREATE TABLE user_email_change_requests (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	email CITEXT NOT NULL,
	status user_email_change_request_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	expires_at TIMESTAMPTZ NOT NULL,
	verified_at TIMESTAMPTZ,
	applied_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT user_email_change_requests_dates CHECK (expires_at > created_at),
	CONSTRAINT user_email_change_requests_verified_time CHECK (
		verified_at IS NULL
		OR verified_at >= created_at
	),
	CONSTRAINT user_email_change_requests_applied_time CHECK (
		applied_at IS NULL
		OR (
			verified_at IS NOT NULL
			AND applied_at >= verified_at
		)
	),
	CONSTRAINT user_email_change_requests_cancelled_time CHECK (
		cancelled_at IS NULL
		OR cancelled_at >= created_at
	)
);

CREATE UNIQUE INDEX user_email_change_requests_one_active_per_user_idx ON user_email_change_requests (user_id)
WHERE status = 'pending';

CREATE UNIQUE INDEX user_email_change_requests_one_active_per_email_idx ON user_email_change_requests (email)
WHERE status = 'pending';

CREATE INDEX user_email_change_requests_expires_idx ON user_email_change_requests (expires_at);
CREATE INDEX user_email_change_requests_user_idx ON user_email_change_requests (user_id);
CREATE INDEX user_email_change_requests_email_idx ON user_email_change_requests (email);

-- Tokens table: store only a hash of the raw token for verification.
CREATE TABLE user_email_change_request_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	user_email_change_request_id UUID NOT NULL REFERENCES user_email_change_requests (id) ON DELETE CASCADE,
	token_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	consumed BOOLEAN NOT NULL DEFAULT false,
	consumed_by UUID REFERENCES users (id) ON DELETE
	SET NULL,
		consumed_at TIMESTAMPTZ,
		CONSTRAINT user_email_change_request_tokens_hash_unique UNIQUE (token_hash),
		CONSTRAINT user_email_change_request_tokens_consumed_time CHECK (
			(consumed = false)
			OR (consumed_at IS NOT NULL)
		)
);

CREATE INDEX user_email_change_request_tokens_hash_idx ON user_email_change_request_tokens (token_hash);
CREATE INDEX user_email_change_request_tokens_request_idx ON user_email_change_request_tokens (user_email_change_request_id);
CREATE TABLE taxa (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	gbif_id INTEGER UNIQUE,
	name CITEXT NOT NULL CONSTRAINT taxon_name_length CHECK (char_length(name) >= 2),
	scientific_name CITEXT GENERATED ALWAYS AS (name || COALESCE(' ' || authorship, '')) STORED,
	rank taxon_rank NOT NULL,
	status taxon_status NOT NULL,
	authorship CITEXT,
	accepted_taxon_id UUID REFERENCES taxa (id) ON DELETE
	SET NULL,
		parent_id UUID REFERENCES taxa (id) ON DELETE CASCADE,
		CONSTRAINT taxon_parent_required_for_non_kingdom CHECK (
			rank = 'Kingdom'
			OR parent_id IS NOT NULL
		),
		CONSTRAINT taxon_synonym_requires_accepted_taxon CHECK (
			(
				status = 'Synonym'
				AND accepted_taxon_id IS NOT NULL
			)
			OR (
				status <> 'Synonym'
				AND accepted_taxon_id IS NULL
			)
		),
		comments TEXT
);
CREATE INDEX taxa_name_rank_status_idx ON taxa (name, rank, status);
CREATE UNIQUE INDEX taxa_name_authorship_uidx ON taxa (name, COALESCE(authorship, ''));
CREATE INDEX taxa_parent_id_idx ON taxa (parent_id);
CREATE INDEX taxa_accepted_taxon_id_idx ON taxa (accepted_taxon_id);
CREATE INDEX taxa_name_trgm_idx ON taxa USING gin (name gin_trgm_ops);
CREATE INDEX taxa_scientific_trgm_idx ON taxa USING gin (scientific_name gin_trgm_ops);
CREATE TABLE taxa_closure (
	ancestor_id UUID NOT NULL REFERENCES taxa(id) ON DELETE CASCADE,
	descendant_id UUID NOT NULL REFERENCES taxa(id) ON DELETE CASCADE,
	depth INTEGER NOT NULL,
	-- 0 = self, 1 = parent, etc.
	PRIMARY KEY (ancestor_id, descendant_id)
);

CREATE INDEX taxa_closure_ancestor_idx ON taxa_closure (ancestor_id);

CREATE INDEX taxa_closure_descendant_idx ON taxa_closure (descendant_id);

CREATE INDEX taxa_closure_depth_idx ON taxa_closure (depth);

CREATE INDEX taxa_closure_cycle_check ON taxa_closure (ancestor_id, descendant_id);

---------------------------------------------------------------------
-- Insert trigger to maintain the closure table
---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION taxa_closure_insert_fn() RETURNS TRIGGER AS $$ BEGIN --
-- 1. self relation
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
VALUES (NEW.id, NEW.id, 0);

-- 2. inherit ancestors from parent
IF NEW.parent_id IS NOT NULL THEN
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
SELECT ancestor_id,
	NEW.id,
	depth + 1
FROM taxa_closure
WHERE descendant_id = NEW.parent_id;
END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_taxa_closure_insert
AFTER
INSERT ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_closure_insert_fn();

---------------------------------------------------------------------
-- Delete trigger to maintain the closure table
---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION taxa_closure_delete_fn() RETURNS TRIGGER AS $$ BEGIN
DELETE FROM taxa_closure
WHERE descendant_id IN (
		SELECT descendant_id
		FROM taxa_closure
		WHERE ancestor_id = OLD.id
	);

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_taxa_closure_delete
AFTER DELETE ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_closure_delete_fn();

---------------------------------------------------------------------
-- Update trigger to maintain the closure table
---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION taxa_closure_update_fn() RETURNS TRIGGER AS $$ BEGIN -- si parent inchangé, rien à faire
	IF NEW.parent_id IS NOT DISTINCT
FROM OLD.parent_id THEN RETURN NEW;
END IF;

-- delete old paths
DELETE FROM taxa_closure
WHERE descendant_id IN (
		SELECT descendant_id
		FROM taxa_closure
		WHERE ancestor_id = OLD.id
	);

-- rebuild self link
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
VALUES (NEW.id, NEW.id, 0);

-- rebuild paths from new parent
IF NEW.parent_id IS NOT NULL THEN
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
SELECT ancestor_id,
	child.descendant_id,
	parent.depth + child.depth + 1
FROM taxa_closure parent
	JOIN taxa_closure child ON child.ancestor_id = OLD.id
WHERE parent.descendant_id = NEW.parent_id;
END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_taxa_closure_update
AFTER
UPDATE OF parent_id ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_closure_update_fn();
-- Prevent taxon <-> parent cycles
CREATE OR REPLACE FUNCTION check_taxa_no_cycle () RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN -- no parent => ok
	IF NEW.parent_id IS NULL THEN RETURN NEW;
END IF;

-- self-parenting
IF NEW.parent_id = NEW.id THEN RAISE EXCEPTION 'Taxon cannot be parent of itself' USING ERRCODE = 'TAX01',
DETAIL = format(
	'Taxon %s cannot be parent of itself',
	NEW.scientific_name
);
END IF;

-- synonym self-check 
IF NEW.accepted_taxon_id IS NOT NULL
AND NEW.accepted_taxon_id = NEW.id THEN RAISE EXCEPTION 'Taxon cannot be synonym of itself' USING ERRCODE = 'TAX02',
DETAIL = format(
	'Taxon %s cannot be synonym of itself',
	NEW.scientific_name
);
END IF;

-- CYCLE CHECK via closure table
IF EXISTS (
	SELECT 1
	FROM taxa_closure
	WHERE ancestor_id = NEW.id
		AND descendant_id = NEW.parent_id
) THEN RAISE EXCEPTION 'Cycle detected in taxon hierarchy' USING ERRCODE = 'TAX01';
END IF;

    RETURN NEW;
END;
$$;
CREATE TRIGGER taxa_no_cycle_trigger BEFORE
INSERT
	OR
UPDATE ON taxa FOR EACH ROW EXECUTE FUNCTION check_taxa_no_cycle ();
-- View to easily query synonyms and their accepted taxa
CREATE VIEW taxon_synonyms AS
SELECT accepted.id AS accepted_taxon_id,
	accepted.name AS accepted_taxon_name,
	synonym.id AS synonym_taxon_id,
	synonym.name AS synonym_taxon_name
FROM taxa accepted
	JOIN taxa synonym ON synonym.accepted_taxon_id = accepted.id;

-- Taxonomic hierarchy view for efficient querying of ancestor-descendant relationships
CREATE VIEW taxon_hierarchy AS WITH RECURSIVE tree AS (
	SELECT t.id,
		t.parent_id,
		t.name,
		t.authorship,
		t.rank,
		t.status,
		t.comments,
		t.id AS root_id,
		ARRAY [t.id] AS path,
		0 AS depth
	FROM taxa t
	WHERE t.parent_id IS NULL
	UNION ALL
	SELECT c.id,
		c.parent_id,
		c.name,
		c.authorship,
		c.rank,
		c.status,
		c.comments,
		tree.root_id,
		tree.path || c.id,
		tree.depth + 1
	FROM taxa c
		JOIN tree ON c.parent_id = tree.id
)
SELECT *
FROM tree;
CREATE TABLE samplings (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	sampling_hash TEXT NOT NULL UNIQUE,
	notes TEXT,
	-- LOCATION FIELDS
	site_code TEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
	site_name TEXT,
	site_locality TEXT,
	site_country_code CHAR(3) REFERENCES countries (code),
	coordinates_precision coordinates_precision,
	coordinates geometry (Point, 4326) NOT NULL,
	latitude REAL GENERATED ALWAYS AS (ST_Y (coordinates)) STORED,
	longitude REAL GENERATED ALWAYS AS (ST_X (coordinates)) STORED,
	altitude INTEGER,
	-- DATE FIELDS
	event_date DATE,
	event_date_precision event_date_precision,
	-- METADATA FIELDS
	performed_by TEXT [],
	-- methods TEXT[] -- FK to events.sampling_method(code)
	-- fixatives TEXT[] -- FK to samples.fixative(code)
	duration INTEGER,
	access_points TEXT [],
	-- R ~ 500m
	h3_res8 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 8)
	) STORED,
	-- R ~ 1km
	h3_res7 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 7)
	) STORED,
	-- R ~ 3 km
	h3_res6 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 6)
	) STORED,
	-- R ~ 10km
	h3_res5 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 5)
	) STORED,
	-- R ~ 25km
	h3_res4 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 4)
	) STORED,
	-- R ~ 60km
	h3_res3 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 3)
	) STORED,
	-- R ~ 166km
	h3_res2 h3index GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 2)
	) STORED
);

CREATE INDEX idx_samplings_hash ON samplings(sampling_hash);
CREATE INDEX samplings_coordinates_gist_idx ON samplings USING GIST (coordinates);
CREATE INDEX samplings_h3_r8_idx ON samplings (h3_res8);
CREATE INDEX samplings_h3_r7_idx ON samplings (h3_res7);
CREATE INDEX samplings_h3_r6_idx ON samplings (h3_res6);
CREATE INDEX samplings_h3_r5_idx ON samplings (h3_res5);
CREATE INDEX samplings_h3_r4_idx ON samplings (h3_res4);
CREATE INDEX samplings_h3_r3_idx ON samplings (h3_res3);
CREATE INDEX samplings_h3_r2_idx ON samplings (h3_res2);

-- Association table linking samplings to taxa targets
CREATE TABLE sampling_target_taxa (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	taxon_id UUID NOT NULL REFERENCES taxa (id) ON DELETE RESTRICT,
	PRIMARY KEY (sampling_id, taxon_id)
);

CREATE INDEX sampling_target_taxa_taxon_idx ON sampling_target_taxa (taxon_id);
-- Habitat groups and habitat tags
CREATE TABLE habitat_groups (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	label CITEXT NOT NULL,
	CONSTRAINT habitat_group_label_unique UNIQUE (label),
	description TEXT,
	exclusive_elements BOOLEAN NOT NULL DEFAULT true
);

-- Habitats belong to a group and may form a hierarchy (parent)
CREATE TABLE habitats (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	label CITEXT NOT NULL,
	description TEXT,
	habitat_group_id UUID NOT NULL REFERENCES habitat_groups (id) ON DELETE CASCADE,
	CONSTRAINT habitat_label_not_empty CHECK (btrim(label) <> ''),
	CONSTRAINT habitat_description_length CHECK (char_length(coalesce(description, '')) <= 4000),
	CONSTRAINT uq_habitat_group_label UNIQUE (habitat_group_id, label)
);

ALTER TABLE habitat_groups
ADD COLUMN parent_habitat_id UUID REFERENCES habitats (id) ON DELETE
SET NULL;


-- Full-text search document column (label + description)
-- ALTER TABLE habitats
-- ADD COLUMN document tsvector GENERATED ALWAYS AS (
-- 		to_tsvector(
-- 			'simple',
-- 			coalesce(label, '') || ' ' || coalesce(description, '')
-- 		)
-- 	) STORED;
-- CREATE INDEX habitat_document_idx ON habitats USING GIN (document);
CREATE INDEX habitat_group_idx ON habitats (habitat_group_id);


-- Association table linking samplings to habitats
CREATE TABLE samplings_habitats (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	habitat_id UUID NOT NULL REFERENCES habitats (id) ON DELETE CASCADE,
	PRIMARY KEY (sampling_id, habitat_id)
);
CREATE INDEX samplings_habitats_habitat_idx ON samplings_habitats (habitat_id);


--
-- Prevent cycles in parent relationship
--
CREATE OR REPLACE FUNCTION validate_group_parent_not_in_subtree()
RETURNS trigger AS $$
DECLARE is_invalid BOOLEAN;
BEGIN IF NEW.parent_habitat_id IS NULL THEN RETURN NEW;
END IF;

    WITH RECURSIVE descendants AS (
	-- tous les habitats du groupe
	SELECT h.id
	FROM habitats h
	WHERE h.habitat_group_id = NEW.id
	UNION ALL
	-- descendance récursive
	SELECT child.id
	FROM habitats child
		JOIN descendants d ON child.parent_id = d.id
)
SELECT TRUE INTO is_invalid
FROM descendants
WHERE id = NEW.parent_habitat_id
LIMIT 1;

    IF is_invalid THEN RAISE EXCEPTION 'Invalid parent_habitat_id: cannot reference a descendant habitat of the same group' USING ERRCODE = 'HB001';
END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TABLE IF NOT EXISTS sampling_methods (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	code CITEXT NOT NULL UNIQUE,
	name CITEXT NOT NULL UNIQUE,
	description TEXT
);


CREATE TABLE IF NOT EXISTS events_sampling_methods (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	method_id UUID NOT NULL REFERENCES sampling_methods (id) ON DELETE CASCADE,
	PRIMARY KEY (sampling_id, method_id)
);

CREATE OR REPLACE FUNCTION resolve_method(
		input_text text,
		min_auto float DEFAULT 0.60
	) RETURNS SETOF vocab_resolution LANGUAGE plpgsql AS $$
DECLARE v_clean text;
DECLARE v_result vocab_resolution;
BEGIN v_clean := lower(trim(input_text));

    ------------------------------------------------------------------
-- EXACT MATCH
------------------------------------------------------------------
SELECT id,
	code,
	name,
	1.0 INTO v_result
FROM sampling_methods
WHERE lower(code) = v_clean
	OR lower(name) = v_clean
LIMIT 1;

    IF FOUND THEN RETURN NEXT vocab_decision(
	v_result.id,
	v_result.code,
	v_result.label,
	1.0,
	min_auto
);
RETURN;
END IF;

------------------------------------------------------------------
-- FUZZY MATCH
------------------------------------------------------------------
SELECT m.id,
	m.code,
	m.name,
	similarity(m.name, v_clean) INTO v_result
FROM sampling_methods m
WHERE m.name % v_clean
ORDER BY similarity(m.name, v_clean) DESC
LIMIT 1;

    IF FOUND THEN RETURN NEXT vocab_decision(
	v_result.id,
	v_result.code,
	v_result.label,
	v_result.confidence,
	min_auto
);
RETURN;
END IF;

    ------------------------------------------------------------------
-- FALLBACK RAW
------------------------------------------------------------------
RETURN NEXT (
	NULL::uuid,
	NULL::text,
	input_text,
	0.0,
	'none',
	false
);

END;
$$;
CREATE TABLE IF NOT EXISTS fixatives (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	code CITEXT NOT NULL UNIQUE,
	name CITEXT NOT NULL UNIQUE,
	description TEXT
);


CREATE TABLE IF NOT EXISTS samplings_fixatives (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	fixative_id UUID NOT NULL REFERENCES fixatives (id) ON DELETE CASCADE,
	PRIMARY KEY (sampling_id, fixative_id)
);

CREATE OR REPLACE FUNCTION resolve_fixative(
		input_text text,
		min_auto float DEFAULT 0.60
	) RETURNS SETOF vocab_resolution LANGUAGE plpgsql AS $$
DECLARE v_clean text;
DECLARE v_result vocab_resolution;
BEGIN v_clean := lower(trim(input_text));

    ------------------------------------------------------------------
-- EXACT MATCH
------------------------------------------------------------------
SELECT id,
	code,
	name,
	1.0 INTO v_result
FROM fixatives
WHERE lower(code) = v_clean
	OR lower(name) = v_clean
LIMIT 1;

    IF FOUND THEN RETURN NEXT vocab_decision(
	v_result.id,
	v_result.code,
	v_result.label,
	1.0,
	min_auto
);
RETURN;
END IF;

    ------------------------------------------------------------------
-- FUZZY MATCH
------------------------------------------------------------------
SELECT m.id,
	m.code,
	m.name,
	similarity(m.name, v_clean) INTO v_result
FROM fixatives m
WHERE m.name % v_clean
ORDER BY similarity(m.name, v_clean) DESC
LIMIT 1;

    IF FOUND THEN RETURN NEXT vocab_decision(
	v_result.id,
	v_result.code,
	v_result.label,
	v_result.confidence,
	min_auto
);
RETURN;
END IF;

    ------------------------------------------------------------------
-- FALLBACK RAW
------------------------------------------------------------------
RETURN NEXT (
	NULL::uuid,
	NULL::text,
	input_text,
	0.0,
	'none',
	false
);

END;
$$;
CREATE TABLE occurrences (
	-- ULID
	id CHAR(26) PRIMARY KEY,
	code TEXT NOT NULL,
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	type_status occurrence_type_status,
	comments TEXT,
	-- Identification fields
	taxon_id UUID NOT NULL REFERENCES taxa (id) ON DELETE RESTRICT,
	verbatim_identification TEXT,
	identified_by TEXT [],
	identification_date DATE,
	identification_date_precision event_date_precision,
	-- whether the identification is a confer (i.e. tentative) identification
	identification_confer BOOLEAN NOT NULL DEFAULT FALSE,
	identification_addendum TEXT,
	-- Content fields
	content_description TEXT,
	quantity_exact INTEGER,
	quantity_lower INTEGER,
	quantity_upper INTEGER,
	sources TEXT [],
	CONSTRAINT occurrence_quantity_shape CHECK (
		(
			quantity_exact IS NOT NULL
			AND quantity_lower IS NULL
			AND quantity_upper IS NULL
		)
		OR (
			quantity_exact IS NULL
			AND (
				quantity_lower IS NOT NULL
				OR quantity_upper IS NOT NULL
			)
			AND (
				quantity_lower IS NULL
				OR quantity_upper IS NULL
				OR quantity_lower <= quantity_upper
			)
		)
	)
);

CREATE INDEX occurrences_sampling_id_idx ON occurrences (sampling_id);

CREATE INDEX occurrences_type_status_idx ON occurrences (type_status);

CREATE TABLE IF NOT EXISTS occurrence_code_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	occurrence_id CHAR(26) NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	code TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE occurrence_collections (
	occurrence_id CHAR(26) NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	vouchers TEXT [],
	CONSTRAINT occurrence_collections_name_not_empty CHECK (char_length(btrim(name)) >= 2),
	PRIMARY KEY (occurrence_id, name)
);
CREATE INDEX occurrence_collections_name_idx ON occurrence_collections (name);
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
-- Minimal dataset and associations
CREATE TABLE datasets (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	import_id UUID NOT NULL REFERENCES import_workflows (import_id) ON DELETE CASCADE UNIQUE,
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
CREATE TABLE import_samplings_occurrences (
	-- =========================
	-- INGESTION CONTEXT
	-- =========================
	import_id UUID NOT NULL REFERENCES import_workflows (import_id) ON DELETE CASCADE,
	imported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	-- =========================
	-- SAMPLING
	-- =========================
	sampling_hash TEXT NOT NULL,
	sampling_notes TEXT,
	site_code TEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
	site_name TEXT,
	site_locality TEXT,
	site_country_code CHAR(3),
	coordinates_precision coordinates_precision,
	longitude DOUBLE PRECISION NOT NULL,
	latitude DOUBLE PRECISION NOT NULL,
	altitude INTEGER,
	event_date DATE,
	event_date_precision event_date_precision,
	performed_by TEXT [],
	duration INTEGER,
	access_points TEXT [],
	-- =========================
	-- OCCURRENCE IDENTIFICATION
	-- =========================
	occurrence_id CHAR(26) NOT NULL,
	occurrence_code TEXT NOT NULL,
	type_status occurrence_type_status,
	taxon_name TEXT,
	taxon_authorship TEXT,
	taxon_scientific_name CITEXT GENERATED ALWAYS AS (
		taxon_name || COALESCE(' ' || taxon_authorship, '')
	) STORED,
	verbatim_identification TEXT,
	identified_by TEXT [],
	identification_date DATE,
	identification_date_precision event_date_precision,
	identification_confer BOOLEAN NOT NULL DEFAULT FALSE,
	identification_addendum TEXT,
	-- =========================
	-- OCCURRENCE METADATA
	-- =========================
	content_description TEXT,
	quantity_exact INTEGER,
	quantity_lower INTEGER,
	quantity_upper INTEGER,
	sources TEXT [],
	occurrence_comments TEXT,
	-- =========================
	-- CONSTRAINTS
	-- =========================
	CONSTRAINT import_sampling_occurrence_qty_check CHECK (
		(
			quantity_exact IS NOT NULL
			AND quantity_lower IS NULL
			AND quantity_upper IS NULL
		)
		OR (
			quantity_exact IS NULL
			AND (
				quantity_lower IS NOT NULL
				OR quantity_upper IS NOT NULL
			)
			AND (
				quantity_lower IS NULL
				OR quantity_upper IS NULL
				OR quantity_lower <= quantity_upper
			)
		)
	)
);

CREATE INDEX idx_staging_import ON import_samplings_occurrences(import_hash);

CREATE INDEX idx_staging_hash ON import_samplings_occurrences(import_hash, sampling_hash);
CREATE TYPE taxon_match_type AS ENUM ('exact', 'name_auth', 'fuzzy', 'name_only');

CREATE TABLE IF NOT EXISTS taxon_candidates (
	import_id UUID NOT NULL REFERENCES import_workflows (import_id) ON DELETE CASCADE,
	input_name TEXT NOT NULL,
	taxon_id UUID,
	match_type taxon_match_type,
	score DOUBLE PRECISION NOT NULL,
	is_ambiguous BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT taxon_candidates_unique UNIQUE (import_hash, input_name, taxon_id)
);
CREATE TABLE IF NOT EXISTS gbif_staging (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	dataset_key UUID,
	occurrence_id CHAR(26),
	verbatim_taxon_name TEXT,
	verbatim_scientific_name TEXT,
	verbatim_identification_qualifier TEXT,
	verbatim_identified_by TEXT,
	verbatim_identification_date TEXT,
	verbatim_identification_date_precision event_date_precision,
	taxon_id UUID,
	identification_confer BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);