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