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