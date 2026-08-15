-- Type alias ULID for 26-character ULID strings
CREATE DOMAIN ULID AS TEXT CHECK (VALUE ~ '^[0-9A-HJKMNP-TV-Z]{26}$');

CREATE DOMAIN doi AS citext CHECK (
	lower(VALUE::text) ~ '^10\.\d{4,9}/[-._;()/:a-z0-9]+$'
);

CREATE TYPE coordinates_precision AS ENUM('<100m', '<1km', '<10km', '10-100km');

CREATE TYPE event_date_precision AS ENUM('day', 'month', 'year');

CREATE TYPE duplicate_source AS ENUM('existing', 'staging');

CREATE TYPE taxon_rank AS ENUM(
	'subspecies',
	'species',
	'subgenus',
	'genus',
	'family',
	'order',
	'class',
	'phylum',
	'kingdom'
);

CREATE TYPE taxon_status AS ENUM(
	'accepted',
	'synonym',
	'doubtful',
	'unreferenced',
	'unclassified'
);

CREATE TYPE occurrence_type_status AS ENUM('holotype', 'neotype', 'topotype');

-- Roles for platform users
CREATE TYPE user_role AS ENUM('visitor', 'contributor', 'maintainer', 'admin');


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