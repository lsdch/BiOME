CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS migrations;

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE EXTENSION IF NOT EXISTS h3;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- used for case-insensitive text type and trigram indexes
CREATE EXTENSION IF NOT EXISTS citext;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Type alias ULID for 26-character ULID strings
CREATE DOMAIN ULID AS TEXT CHECK (VALUE ~ '^[0-9A-HJKMNP-TV-Z]{26}$');

CREATE TYPE coordinates_precision AS ENUM('<100m', '<1km', '<10km', '10-100km');

CREATE TYPE event_date_precision AS ENUM('day', 'month', 'year');

CREATE TYPE duplicate_source AS ENUM('existing', 'staging');

CREATE TYPE taxon_rank AS ENUM(
	'SUBSPECIES',
	'SPECIES',
	'SUBGENUS',
	'GENUS',
	'FAMILY',
	'ORDER',
	'CLASS',
	'PHYLUM',
	'KINGDOM'
);

CREATE TYPE taxon_status AS ENUM(
	'ACCEPTED',
	'SYNONYM',
	'DOUBTFUL',
	'UNREFERENCED',
	'UNCLASSIFIED'
);

CREATE TYPE occurrence_type_status AS ENUM('HOLOTYPE', 'NEOTYPE', 'TOPOTYPE');

-- Roles for platform users
CREATE TYPE user_role AS ENUM('Visitor', 'Contributor', 'Maintainer', 'Admin');


CREATE TYPE sort_direction AS ENUM('asc', 'desc');

CREATE TYPE occurrence_order_by AS ENUM(
	'code',
	'site_name',
	'site_code',
	'event_date',
	'taxon_name',
	'created_at',
	'updated_at'
);

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

CREATE TABLE IF NOT EXISTS settings (
    -- enforce singleton
    id SERIAL PRIMARY KEY CHECK (id = 1),
    app_name TEXT NOT NULL DEFAULT 'BiOME',
    app_subtitle TEXT,
    app_description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    account_requests_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    admin_email TEXT NOT NULL,
    -- email settings
    smtp_settings_valid BOOLEAN NOT NULL DEFAULT FALSE,
    -- indicates whether SMTP settings are valid (e.g. can connect to SMTP server with provided settings)
    smtp_host TEXT,
    smtp_port INTEGER,
    smtp_user TEXT,
    smtp_password TEXT,
    smtp_from_email TEXT,
    smtp_from_name TEXT,
    -- external services
    geoapify_api_key TEXT,
    geoapify_usage_limit INTEGER NOT NULL DEFAULT 3000,
    -- GBIF contact email for occurrence downloads
    -- is sent along GBIF requests to identify the user and allow GBIF to contact them if needed
    gbif_contact_email TEXT NOT NULL
)

CREATE TABLE IF NOT EXISTS geoapify_usage (
    id SERIAL PRIMARY KEY,
    usage_date DATE NOT NULL,
    requests_count INTEGER NOT NULL,
    UNIQUE (usage_date)
)

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

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

-- Invitations and tokens (separate metadata and token storage)
-- This design separates invitation metadata from token hashes and allows multiple tokens per invitation (rotation).
CREATE TYPE invitation_status AS ENUM('pending', 'redeemed', 'cancelled', 'expired');
CREATE TABLE invitations (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	email CITEXT NOT NULL,
	invitee_name TEXT NOT NULL,
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
--   VALUES (:inv_id, token_sha256(:raw_token));

-- User account requests and email-change requests.
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

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    refresh_token_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    user_agent TEXT,
    ip_address INET
);

CREATE TABLE taxa (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    gbif_id INTEGER UNIQUE,
    name CITEXT NOT NULL CONSTRAINT taxon_name_length CHECK (char_length(name) >= 2),
    scientific_name CITEXT NOT NULL GENERATED ALWAYS AS (name || COALESCE(' ' || authorship, '')) STORED,
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
	coordinates_precision INTEGER,
	coordinates geometry (Point, 4326) NOT NULL,
	latitude REAL NOT NULL GENERATED ALWAYS AS (ST_Y (coordinates)) STORED,
	longitude REAL NOT NULL GENERATED ALWAYS AS (ST_X (coordinates)) STORED,
	altitude INTEGER,
	-- DATE FIELDS
	event_date DATE,
	event_date_precision event_date_precision,
	-- METADATA FIELDS
	performed_by TEXT [],
	duration INTEGER,
	access_points TEXT [],
	-- R ~ 500m
	h3_index BIGINT NOT NULL GENERATED ALWAYS AS (
		h3_lat_lng_to_cell(point(ST_X(coordinates), ST_Y(coordinates)), 12)::BIGINT
	) STORED,
	CONSTRAINT samplings_coordinates_precision_range CHECK (
		coordinates_precision IS NULL
		OR (
			coordinates_precision >= 0
			AND coordinates_precision <= 100000
		)
	)
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
	exclusive_elements BOOLEAN NOT NULL DEFAULT true,
	CONSTRAINT habitat_group_label_not_empty CHECK (btrim(label) <> '')
);

CREATE UNIQUE INDEX idx_habitat_group_label_uq ON habitat_groups (label);

-- Habitats belong to a group and may form a hierarchy (parent)
CREATE TABLE habitats (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	label CITEXT NOT NULL,
	description TEXT,
	habitat_group_id UUID NOT NULL REFERENCES habitat_groups (id) ON DELETE CASCADE,
	CONSTRAINT habitat_label_not_empty CHECK (btrim(label) <> ''),
	CONSTRAINT habitat_description_length CHECK (char_length(coalesce(description, '')) <= 4000),
	CONSTRAINT uq_habitat_label UNIQUE (label)
);

ALTER TABLE habitat_groups
ADD COLUMN parent_habitat_id UUID REFERENCES habitats (id) ON DELETE
SET NULL;

CREATE INDEX habitat_group_idx ON habitats (habitat_group_id);

-- Full-text search document column (label + description)
-- ALTER TABLE habitats
-- ADD COLUMN document tsvector GENERATED ALWAYS AS (
-- 		to_tsvector(
-- 			'simple',
-- 			coalesce(label, '') || ' ' || coalesce(description, '')
-- 		)
-- 	) STORED;
-- CREATE INDEX habitat_document_idx ON habitats USING GIN (document);
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

CREATE TABLE IF NOT EXISTS occurrences (
	-- ULID
	id ULID PRIMARY KEY,
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
	-- Metadata fields
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	import_batch_id ULID REFERENCES import_batches (id) ON DELETE
	SET NULL,
		-- Constraints
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

CREATE TRIGGER occurrences_set_updated_at BEFORE
UPDATE ON occurrences FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS occurrence_code_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	code TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE occurrence_collections (
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
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
CREATE TABLE occurrences_articles (
	occurrence_id ULID NOT NULL REFERENCES occurrences (id) ON DELETE CASCADE,
	article_id UUID NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
	PRIMARY KEY (occurrence_id, article_id)
);

CREATE INDEX occurrences_articles_article_idx ON occurrences_articles (article_id);

-- Minimal dataset and associations
CREATE TABLE datasets (
	id ULID PRIMARY KEY DEFAULT gen_random_uuid (),
	label TEXT NOT NULL,
	slug TEXT NOT NULL,
	description TEXT,
	assembled_by TEXT [],
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
CREATE TABLE dataset_curator (
	dataset_id ULID NOT NULL REFERENCES datasets (id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	PRIMARY KEY (dataset_id, user_id)
);

CREATE INDEX dataset_curator_user_idx ON dataset_curator (user_id);

CREATE TABLE IF NOT EXISTS import_batches (
    id ULID PRIMARY KEY,
    label TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    submitted_by TEXT,
    assembled_by TEXT [],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT import_batches_label_length CHECK (
        CHAR_LENGTH(BTRIM(label)) BETWEEN 4 AND 40
    ),
    CONSTRAINT import_batches_slug_length CHECK (
        CHAR_LENGTH(BTRIM(slug)) BETWEEN 4 AND 40
    )
);

CREATE TYPE gbif_import_status AS ENUM (
    'pending',
    'in_progress',
    'completed',
    'failed'
);

CREATE TABLE import_workflows (
    import_hash TEXT PRIMARY KEY,
    label TEXT,
    gbif_status gbif_import_status NOT NULL DEFAULT 'pending',
    gbif_candidates_total INTEGER,
    gbif_candidates_fetched INTEGER,
    gbif_claimed_at TIMESTAMPTZ,
    gbif_updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ
);


CREATE TABLE import_samplings_occurrences (
    -- =========================
    -- INGESTION CONTEXT
    -- =========================
    import_hash TEXT NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    row_number INTEGER NOT NULL,
    PRIMARY KEY (import_hash, row_number),
    -- =========================
    -- SAMPLING
    -- =========================
    sampling_hash TEXT NOT NULL,
    sampling_comments TEXT,
    site_code TEXT CHECK (site_code ~ '^[A-Za-z0-9.-]{3,32}$'),
    site_name TEXT,
    site_locality TEXT,
    site_country_code CHAR(3),
    coordinates_precision INTEGER,
    longitude REAL NOT NULL,
    latitude REAL NOT NULL,
    coordinates geometry (Point, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
    ) STORED,
    altitude INTEGER,
    event_date DATE,
    event_date_precision event_date_precision,
    performed_by TEXT [],
    duration INTEGER,
    access_points CITEXT [],
    sampling_targets CITEXT [],
    sampling_fixatives CITEXT [],
    sampling_methods CITEXT [],
    habitats CITEXT [],
    -- =========================
    -- OCCURRENCE IDENTIFICATION
    -- =========================
    -- occurrence_code can be provided by the user, but will not be used as the actual occurrence code.
    -- Instead, the occurrence_code will be generated.
    -- The provided occurrence_code will be stored as a reference for the user, in the codes history.
    occurrence_code TEXT,
    type_status occurrence_type_status,
    taxon_name CITEXT NOT NULL,
    taxon_authorship CITEXT,
    taxon_scientific_name CITEXT NOT NULL GENERATED ALWAYS AS (
        taxon_name || COALESCE(' ' || taxon_authorship, '')
    ) STORED,
    taxon_rank CITEXT,
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
    CONSTRAINT import_coordinate_precision_check CHECK (
        coordinates_precision IS NULL
        OR (
            coordinates_precision >= 0
            AND coordinates_precision <= 100000
        )
    ),
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

CREATE TABLE IF NOT EXISTS gbif_staging (
    key INTEGER PRIMARY KEY,
    parent TEXT,
    parent_key INTEGER,
    canonical_name TEXT NOT NULL,
    scientific_name TEXT NOT NULL,
    status TEXT NOT NULL,
    rank TEXT NOT NULL,
    name_type TEXT NOT NULL,
    kingdom_key INTEGER,
    phylum_key INTEGER,
    class_key INTEGER,
    order_key INTEGER,
    family_key INTEGER,
    genus_key INTEGER,
    species_key INTEGER,
    higher_taxon_keys INTEGER [],
    higher_taxon_names TEXT [],
    authorship text,
    num_descendants INTEGER,
    accepted_key INTEGER,
    accepted_name TEXT
);

CREATE TABLE IF NOT EXISTS gbif_dependencies (
    import_hash TEXT NOT NULL,
    key INTEGER NOT NULL,
    PRIMARY KEY (import_hash, key)
);

CREATE TYPE taxon_match_type AS ENUM ('exact', 'fuzzy', 'name_only');
CREATE TYPE taxon_match_source AS ENUM ('internal', 'gbif', 'manual');

CREATE TABLE IF NOT EXISTS taxon_candidates (
    import_hash TEXT NOT NULL,
    input_name CITEXT NOT NULL,
    source taxon_match_source NOT NULL,
    match_type taxon_match_type NOT NULL,
    taxon_id UUID REFERENCES taxa(id),
    gbif_id INTEGER REFERENCES gbif_staging(key),
    score DOUBLE PRECISION,
    priority INTEGER NOT NULL,
    name CITEXT NOT NULL,
    authorship CITEXT,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    CONSTRAINT taxon_candidates_target_check CHECK (
        (
            source = 'internal'
            AND taxon_id IS NOT NULL
        )
        OR (
            source = 'gbif'
            AND gbif_id IS NOT NULL
            AND taxon_id IS NULL
        )
    )
);

CREATE UNIQUE INDEX taxon_candidates_internal_unique ON taxon_candidates (import_hash, input_name, taxon_id)
WHERE source = 'internal';

CREATE UNIQUE INDEX taxon_candidates_gbif_unique ON taxon_candidates (import_hash, input_name, gbif_id)
WHERE source = 'gbif';

CREATE TABLE taxa_staging (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    import_hash TEXT NOT NULL,
    name CITEXT NOT NULL,
    authorship TEXT,
    rank taxon_rank NOT NULL,
    status taxon_status NOT NULL,
    parent_source taxon_match_source NOT NULL,
    parent_taxa_id UUID REFERENCES taxa(id),
    parent_gbif_id INTEGER REFERENCES gbif_staging(key),
    parent_input_name TEXT,
    CHECK (
        (
            parent_source = 'internal'
            AND parent_taxa_id IS NOT NULL
        )
        OR (
            parent_source = 'gbif'
            AND parent_gbif_id IS NOT NULL
        )
        OR (
            parent_source = 'manual'
            AND parent_input_name IS NOT NULL
        )
    )
);

CREATE TYPE resolution_status AS ENUM (
    'pending',
    'auto_resolved',
    'user_resolved',
    'needs_decision'
);

CREATE TYPE taxon_gbif_status AS ENUM (
    'skipped',
    'pending',
    'completed',
    'failed'
);

CREATE TABLE IF NOT EXISTS taxon_resolution (
    import_hash TEXT NOT NULL,
    input_name CITEXT NOT NULL,
    source taxon_match_source,
    gbif_id INTEGER REFERENCES gbif_staging(key),
    taxon_id UUID REFERENCES taxa(id),
    staging_id UUID REFERENCES taxa_staging(id),
    status resolution_status DEFAULT 'pending',
    gbif_status taxon_gbif_status DEFAULT 'skipped',
    PRIMARY KEY (import_hash, input_name)
);

CREATE TYPE vocab_resolution_status AS ENUM (
    'auto',
    'selected',
    'pending',
    'request_creation',
    'discard'
);


CREATE TABLE IF NOT EXISTS sampling_methods_resolution (
    import_hash TEXT NOT NULL,
    input_text TEXT NOT NULL,
    resolved_method_id UUID REFERENCES sampling_methods (id) ON DELETE CASCADE,
    status vocab_resolution_status NOT NULL DEFAULT 'pending',
    PRIMARY KEY (import_hash, input_text)
)

