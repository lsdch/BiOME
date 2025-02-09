CREATE MIGRATION m1uf7gqulw57cndzcygywqjen4geft6x3e7ckzsrfk6vhqzby46aaq
    ONTO initial
{
  CREATE EXTENSION pgcrypto VERSION '1.3';
  CREATE EXTENSION auth VERSION '1.0';
  CREATE EXTENSION pg_trgm VERSION '1.6';
  CREATE MODULE admin IF NOT EXISTS;
  CREATE MODULE datasets IF NOT EXISTS;
  CREATE MODULE date IF NOT EXISTS;
  CREATE MODULE events IF NOT EXISTS;
  CREATE MODULE location IF NOT EXISTS;
  CREATE MODULE occurrence IF NOT EXISTS;
  CREATE MODULE people IF NOT EXISTS;
  CREATE MODULE references IF NOT EXISTS;
  CREATE MODULE samples IF NOT EXISTS;
  CREATE MODULE sampling IF NOT EXISTS;
  CREATE MODULE seq IF NOT EXISTS;
  CREATE MODULE storage IF NOT EXISTS;
  CREATE MODULE taxonomy IF NOT EXISTS;
  CREATE MODULE tokens IF NOT EXISTS;
  CREATE MODULE traits IF NOT EXISTS;
  CREATE GLOBAL default::current_user_id -> std::uuid {
      SET default := (<std::uuid>{});
  };
  CREATE TYPE location::Country {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(2);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE ANNOTATION std::description := 'Countries as defined in the ISO 3166-1 norm.';
  };
  CREATE SCALAR TYPE location::CoordinatesPrecision EXTENDING enum<`<100m`, `<1km`, `<10km`, `10-100km`, Unknown>;
  CREATE ABSTRACT TYPE default::Auditable {
      CREATE ANNOTATION std::title := 'Auto-generation of timestamps';
  };
  CREATE TYPE location::Site EXTENDING default::Auditable {
      CREATE REQUIRED LINK country: location::Country;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(10);
          CREATE CONSTRAINT std::min_len_value(3);
          CREATE ANNOTATION std::description := 'A short, unique, user-generated, alphanumeric identifier. Recommended size is 8.';
          CREATE ANNOTATION std::title := 'Site identifier';
      };
      CREATE INDEX ON (.code);
      CREATE PROPERTY altitude: std::int32 {
          CREATE ANNOTATION std::title := 'The site elevation in meters';
      };
      CREATE REQUIRED PROPERTY coordinates: tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32> {
          CREATE CONSTRAINT std::expression ON (((((.latitude <= 90) AND (.latitude >= -90)) AND (.longitude <= 180)) AND (.longitude >= -180)));
          CREATE REWRITE
              INSERT 
              USING ((
                  precision := __subject__.coordinates.precision,
                  latitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.latitude, 5),
                  longitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.longitude, 5)
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  precision := __subject__.coordinates.precision,
                  latitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.latitude, 5),
                  longitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.longitude, 5)
              ));
      };
      CREATE PROPERTY description: std::str;
      CREATE PROPERTY locality: std::str;
      CREATE REQUIRED PROPERTY name: std::str;
      CREATE REQUIRED PROPERTY user_defined_locality: std::bool {
          SET default := false;
      };
  };
  ALTER TYPE location::Country {
      CREATE MULTI LINK sites := (.<country[IS location::Site]);
  };
  CREATE ALIAS location::CountryList := (
      SELECT
          location::Country {
              *,
              sites_count := std::count(.sites)
          }
  );
  CREATE ABSTRACT TYPE default::CodeIdentifier {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY code_history: array<tuple<code: std::str, time: std::datetime>> {
          CREATE REWRITE
              UPDATE 
              USING ((IF (__old__.code != .code) THEN ((__old__.code_history ?? []) ++ [(
                  code := __old__.code,
                  time := std::datetime_of_statement()
              )]) ELSE .code_history));
      };
  };
  CREATE FUNCTION references::generate_article_code(authors: array<std::str>, year: std::int32) ->  std::str USING (SELECT
      (IF (std::len(authors) = 1) THEN (((std::str_split((authors)[0], ' '))[0] ++ '_') ++ <std::str>year) ELSE (IF (std::len(authors) = 2) THEN ((std::array_join([(std::str_split((authors)[0], ' '))[0], (std::str_split((authors)[1], ' '))[0]], '_') ++ '_') ++ <std::str>year) ELSE (((std::str_split((authors)[0], ' '))[0] ++ '_et_al_') ++ <std::str>year)))
  );
  CREATE SCALAR TYPE occurrence::OccurrenceCategory EXTENDING enum<Internal, External>;
  CREATE ABSTRACT TYPE occurrence::Occurrence EXTENDING default::Auditable {
      CREATE PROPERTY comments: std::str;
  };
  CREATE ABSTRACT TYPE occurrence::BioMaterial EXTENDING default::CodeIdentifier, occurrence::Occurrence {
      ALTER PROPERTY code {
          SET OWNED;
          SET REQUIRED;
          SET TYPE std::str;
          ALTER CONSTRAINT std::exclusive {
              SET OWNED;
          };
          CREATE ANNOTATION std::description := "Format like 'taxon_short_code[sampling_code]'";
      };
      CREATE REQUIRED PROPERTY is_type: std::bool {
          SET default := false;
      };
  };
  CREATE SCALAR TYPE date::DatePrecision EXTENDING enum<Year, Month, Day, Unknown>;
  CREATE TYPE occurrence::Identification EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY identified_on: tuple<date: std::datetime, precision: date::DatePrecision>;
  };
  ALTER TYPE occurrence::Occurrence {
      CREATE REQUIRED LINK identification: occurrence::Identification {
          ON SOURCE DELETE DELETE TARGET;
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE SCALAR TYPE occurrence::QuantityType EXTENDING enum<Unknown, One, Several, Ten, Tens, Hundred>;
  CREATE TYPE occurrence::ExternalBioMat EXTENDING occurrence::BioMaterial {
      CREATE PROPERTY content_description: std::str;
      CREATE PROPERTY in_collection: std::str;
      CREATE MULTI PROPERTY item_vouchers: std::str;
      CREATE PROPERTY original_link: std::str;
      CREATE REQUIRED PROPERTY quantity: occurrence::QuantityType;
  };
  CREATE ABSTRACT TYPE default::Vocabulary {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE DELEGATED CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(12);
          CREATE CONSTRAINT std::min_len_value(2);
          CREATE ANNOTATION std::title := 'An expressive, unique, user-generated uppercase alphanumeric code';
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE INDEX ON ((.code, .label));
      CREATE ANNOTATION std::title := 'An extensible list of terms';
  };
  CREATE SCALAR TYPE seq::ExtSeqOrigin EXTENDING enum<Lab, PersCom, DB>;
  CREATE ABSTRACT TYPE seq::Sequence EXTENDING default::Auditable, default::CodeIdentifier {
      CREATE REQUIRED LINK identification: occurrence::Identification;
      CREATE PROPERTY comments: std::str;
      CREATE PROPERTY label: std::str;
      CREATE PROPERTY legacy: tuple<id: std::int32, code: std::str, alignment_code: std::str> {
          CREATE ANNOTATION std::description := 'Legacy identifiers for retrocompatibility with data stored in GOTIT.';
      };
      CREATE PROPERTY sequence: std::str;
  };
  CREATE TYPE seq::ExternalSequence EXTENDING seq::Sequence, occurrence::Occurrence {
      CREATE LINK source_sample: occurrence::ExternalBioMat;
      ALTER PROPERTY code {
          SET OWNED;
          SET REQUIRED;
          SET TYPE std::str;
          ALTER CONSTRAINT std::exclusive {
              SET OWNED;
          };
      };
      CREATE REQUIRED PROPERTY origin: seq::ExtSeqOrigin;
      CREATE REQUIRED PROPERTY specimen_identifier: std::str {
          CREATE ANNOTATION std::description := 'An identifier for the organism from which the sequence was produced, provided in the original source';
      };
      CREATE PROPERTY original_taxon: std::str {
          CREATE ANNOTATION std::description := 'The verbatim identification provided in the original source.';
      };
  };
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE MULTI LINK sequences := (.<source_sample[IS seq::ExternalSequence]);
  };
  CREATE SCALAR TYPE taxonomy::Rank EXTENDING enum<Kingdom, Phylum, Class, Order, Family, Genus, Species, Subspecies>;
  CREATE SCALAR TYPE taxonomy::TaxonStatus EXTENDING enum<Accepted, Unreferenced, Unclassified>;
  CREATE TYPE taxonomy::Taxon EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(4);
      };
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE REWRITE
              INSERT 
              USING ((IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE std::str_replace(.name, ' ', '_')));
          CREATE REWRITE
              UPDATE 
              USING ((IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE std::str_replace(.name, ' ', '_')));
      };
      CREATE REQUIRED PROPERTY status: taxonomy::TaxonStatus;
      CREATE CONSTRAINT std::exclusive ON ((.name, .status));
      CREATE REQUIRED PROPERTY rank: taxonomy::Rank;
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) >= 3)) EXCEPT ((.rank != taxonomy::Rank.Subspecies)) {
          SET errmessage := 'A subspecies name must include at least 2 whitespaces.';
      };
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) >= 2)) EXCEPT ((.rank != taxonomy::Rank.Species)) {
          SET errmessage := 'A species name must include a whitespace.';
      };
      CREATE CONSTRAINT std::expression ON (NOT (std::contains(.name, ' '))) EXCEPT (((.rank = taxonomy::Rank.Species) OR (.rank = taxonomy::Rank.Subspecies))) {
          SET errmessage := 'Taxon names with rank higher than species may not include a whitespace.';
      };
      CREATE LINK parent: taxonomy::Taxon {
          ON TARGET DELETE DELETE SOURCE;
      };
      CREATE CONSTRAINT std::expression ON (EXISTS (.parent)) EXCEPT ((.rank = taxonomy::Rank.Kingdom));
      CREATE INDEX ON ((.name, .code, .rank, .status));
      CREATE MULTI LINK children := (.<parent[IS taxonomy::Taxon]);
      CREATE REQUIRED PROPERTY children_count := (SELECT
          std::count(.children)
      );
      CREATE LINK class: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Class) THEN .parent ELSE .parent.class) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Class) THEN .parent ELSE .parent.class) ELSE <taxonomy::Taxon>{}));
      };
      CREATE LINK family: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Family) THEN .parent ELSE .parent.family) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Family) THEN .parent ELSE .parent.family) ELSE <taxonomy::Taxon>{}));
      };
      CREATE LINK genus: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Genus) THEN .parent ELSE .parent.genus) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Genus) THEN .parent ELSE .parent.genus) ELSE <taxonomy::Taxon>{}));
      };
      CREATE LINK kingdom: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Kingdom) THEN .parent ELSE .parent.kingdom) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Kingdom) THEN .parent ELSE .parent.kingdom) ELSE <taxonomy::Taxon>{}));
      };
      CREATE LINK order: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Order) THEN .parent ELSE .parent.order) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Order) THEN .parent ELSE .parent.order) ELSE <taxonomy::Taxon>{}));
      };
      CREATE LINK phylum: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Phylum) THEN .parent ELSE .parent.phylum) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Phylum) THEN .parent ELSE .parent.phylum) ELSE <taxonomy::Taxon>{}));
      };
      CREATE LINK species: taxonomy::Taxon {
          ON TARGET DELETE ALLOW;
          CREATE REWRITE
              INSERT 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Species) THEN .parent ELSE .parent.species) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Species) THEN .parent ELSE .parent.species) ELSE <taxonomy::Taxon>{}));
      };
      CREATE PROPERTY GBIF_ID: std::int32 {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY anchor: std::bool {
          SET default := false;
          CREATE ANNOTATION std::description := 'Signals whether this taxon was manually imported';
      };
      CREATE PROPERTY authorship: std::str;
      CREATE PROPERTY comment: std::str;
  };
  ALTER TYPE occurrence::Identification {
      CREATE REQUIRED LINK taxon: taxonomy::Taxon;
  };
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE REQUIRED PROPERTY is_homogenous := (SELECT
          (std::count(DISTINCT (.sequences.identification.taxon)) <= 1)
      );
      CREATE SINGLE LINK seq_consensus := (SELECT
          (IF .is_homogenous THEN std::assert_single(DISTINCT (.sequences.identification.taxon), message := ('BioMaterial is marked as homogenous, yet specimens have identification mismatch. UUID: ' ++ <std::str>.id)) ELSE {})
      );
      CREATE REQUIRED PROPERTY is_congruent := (SELECT
          std::assert_exists((.is_homogenous AND (NOT (EXISTS (.sequences)) OR (.identification.taxon IN std::assert_single(DISTINCT (.sequences.identification.taxon), message := ('BioMaterial is marked as homogenous, yet specimens have identification mismatch. UUID: ' ++ <std::str>.id))))))
      );
  };
  CREATE TYPE occurrence::InternalBioMat EXTENDING occurrence::BioMaterial;
  CREATE ABSTRACT TYPE samples::Sample EXTENDING default::Auditable {
      CREATE REQUIRED LINK biomat: occurrence::InternalBioMat;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY number: std::int16 {
          CREATE ANNOTATION std::description := 'Incremental number that discriminates between tubes having the same type in a bio material lot. Used to generate the tube code.';
      };
  };
  CREATE TYPE default::Picture {
      CREATE PROPERTY legend: std::str;
      CREATE REQUIRED PROPERTY path: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE samples::Specimen EXTENDING samples::Sample {
      CREATE MULTI LINK pictures: default::Picture {
          ON SOURCE DELETE DELETE TARGET;
      };
      CREATE REQUIRED PROPERTY molecular_number: std::str;
      CREATE REQUIRED PROPERTY morphological_code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE ANNOTATION std::description := 'A single specimen isolated in a tube.';
      CREATE PROPERTY molecular_code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE MULTI LINK specimens := (.<biomat[IS samples::Specimen]);
  };
  CREATE TYPE seq::AssembledSequence EXTENDING seq::Sequence {
      CREATE REQUIRED LINK specimen: samples::Specimen;
      ALTER LINK identification {
          ON SOURCE DELETE DELETE TARGET;
          SET OWNED;
          SET REQUIRED;
          SET TYPE occurrence::Identification;
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY is_reference: std::bool {
          CREATE ANNOTATION std::description := 'Whether this sequence should be used as the reference for the identification of the specimen';
      };
      CREATE CONSTRAINT std::exclusive ON ((.specimen, .is_reference)) EXCEPT (NOT (.is_reference));
      CREATE PROPERTY accession_number: std::str {
          CREATE ANNOTATION std::description := 'The NCBI accession number, if the sequence was uploaded.';
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY alignmentCode: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE samples::Specimen {
      CREATE MULTI LINK sequences := (.<specimen[IS seq::AssembledSequence]);
      CREATE LINK identification := (((SELECT
          .sequences
      FILTER
          .is_reference
      )).identification);
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE REQUIRED PROPERTY is_homogenous := (SELECT
          (std::count(DISTINCT (.specimens.identification.taxon)) <= 1)
      );
      CREATE SINGLE LINK seq_consensus := (SELECT
          (IF .is_homogenous THEN std::assert_single(DISTINCT (.specimens.identification.taxon), message := ('BioMaterial is marked as homogenous, yet specimens have identification mismatch. UUID: ' ++ <std::str>.id)) ELSE {})
      );
      CREATE REQUIRED PROPERTY is_congruent := (SELECT
          std::assert_exists((.is_homogenous AND (NOT (EXISTS (.specimens)) OR (.identification.taxon IN std::assert_single(DISTINCT (.specimens.identification.taxon), message := ('BioMaterial is marked as homogenous, yet specimens have identification mismatch. UUID: ' ++ <std::str>.id))))))
      );
      CREATE MULTI LINK identified_taxa := (SELECT
          (DISTINCT (.specimens.identification.taxon) ?? .identification.taxon)
      );
  };
  ALTER TYPE occurrence::BioMaterial {
      CREATE REQUIRED PROPERTY category := (std::assert_exists((IF (__source__ IS occurrence::InternalBioMat) THEN occurrence::OccurrenceCategory.Internal ELSE (IF (__source__ IS occurrence::ExternalBioMat) THEN occurrence::OccurrenceCategory.External ELSE <occurrence::OccurrenceCategory>{})), message := (('Occurrence category for occurrence::BioMaterial subtype ' ++ __source__.__type__.name) ++ ' is undefined')));
  };
  CREATE ALIAS occurrence::BioMaterialWithType := (
      SELECT
          occurrence::BioMaterial {
              required has_sequences := EXISTS (([IS occurrence::ExternalBioMat].sequences ?? [IS occurrence::InternalBioMat].specimens.sequences)),
              required is_homogenous := (([IS occurrence::ExternalBioMat].is_homogenous ?? true) AND ([IS occurrence::ExternalBioMat].is_homogenous ?? true)),
              required is_congruent := (<std::bool>[IS occurrence::ExternalBioMat].is_congruent ?? (<std::bool>[IS occurrence::InternalBioMat].is_congruent ?? true)),
              seq_consensus := (<taxonomy::Taxon>[IS occurrence::ExternalBioMat].seq_consensus ?? <taxonomy::Taxon>[IS occurrence::InternalBioMat].seq_consensus),
              external := [IS occurrence::ExternalBioMat] {
                  original_link,
                  in_collection,
                  item_vouchers,
                  quantity,
                  content_description
              }
          }
  );
  CREATE ABSTRACT ANNOTATION default::example;
  CREATE SCALAR TYPE events::SamplingNumber EXTENDING std::sequence;
  CREATE SCALAR TYPE events::SamplingTarget EXTENDING enum<Community, Unknown, Taxa>;
  CREATE SCALAR TYPE people::OrgKind EXTENDING enum<Lab, FundingAgency, SequencingPlatform, Other>;
  CREATE SCALAR TYPE people::UserRole EXTENDING enum<Visitor, Contributor, Maintainer, Admin>;
  CREATE SCALAR TYPE seq::ChromatoQuality EXTENDING enum<Contaminated, Failure, Ok, Unknown>;
  CREATE SCALAR TYPE seq::DNAQuality EXTENDING enum<Unknown, Contaminated, Bad, Good>;
  CREATE SCALAR TYPE seq::PCRQuality EXTENDING enum<Failure, Acceptable, Good, Unknown>;
  CREATE SCALAR TYPE traits::Category EXTENDING enum<Morphology, Physiology, Ecology, Behaviour, LifeHistory, HabitatPref>;
  CREATE SCALAR TYPE traits::TraitDefinitionScope EXTENDING enum<Specimen, Taxon>;
  CREATE ABSTRACT CONSTRAINT date::required_unless_unknown(date: std::datetime, precision: date::DatePrecision) {
      SET errmessage := "Date value is required, except when precision is 'Unknown'";
      USING ((EXISTS (date) IF (precision != date::DatePrecision.Unknown) ELSE NOT (EXISTS (date))));
  };
  CREATE TYPE events::Event EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY performed_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(__subject__.date, __subject__.precision);
      };
      CREATE REQUIRED LINK site: location::Site {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE DELETE SOURCE;
      };
      CREATE REQUIRED PROPERTY code := (WITH
          date := 
              .performed_on.date
          ,
          precision := 
              .performed_on.precision
      SELECT
          ((.site.code ++ '|') ++ (IF (precision = date::DatePrecision.Unknown) THEN 'undated' ELSE (IF (precision = date::DatePrecision.Year) THEN <std::str>std::datetime_get(date, 'year') ELSE ((<std::str>std::datetime_get(date, 'year') ++ '-') ++ <std::str>std::datetime_get(date, 'month')))))
      );
      CREATE MULTI LINK spottings: taxonomy::Taxon;
      CREATE PROPERTY comments: std::str;
  };
  ALTER TYPE occurrence::Identification {
      ALTER PROPERTY identified_on {
          CREATE CONSTRAINT date::required_unless_unknown(__subject__.date, __subject__.precision);
      };
  };
  CREATE TYPE samples::Slide EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY created_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE MULTI LINK pictures: default::Picture {
          ON SOURCE DELETE DELETE TARGET;
      };
      CREATE REQUIRED PROPERTY storage_position: std::int16;
      CREATE REQUIRED LINK specimen: samples::Specimen;
      CREATE PROPERTY code: std::str {
          CREATE ANNOTATION std::description := "Generated as '{collectionCode}_{containerCode}_{slidePositionInBox}'";
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comment: std::str;
      CREATE REQUIRED PROPERTY label: std::str;
  };
  CREATE TYPE seq::DNA EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY extracted_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE REQUIRED LINK specimen: samples::Specimen;
      CREATE REQUIRED PROPERTY chelex_tube: tuple<color: std::str, number: std::int16>;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comments: std::str;
      CREATE PROPERTY concentration: std::float32 {
          CREATE ANNOTATION std::title := 'DNA concentration in ng/µL';
          CREATE CONSTRAINT std::min_value(1e-3);
      };
      CREATE REQUIRED PROPERTY is_empty: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY quality: seq::DNAQuality;
  };
  CREATE TYPE seq::PCR EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY performed_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE REQUIRED LINK DNA: seq::DNA;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY quality: seq::PCRQuality;
  };
  CREATE FUNCTION date::from_json_with_precision(value: std::json) ->  tuple<date: std::datetime, precision: date::DatePrecision> USING (std::assert_exists((
      date := (IF EXISTS ((value)['date']) THEN std::to_datetime(<std::int64>(value)['date']['year'], (<std::int64>std::json_get(value, 'date', 'month') ?? 1), (<std::int64>std::json_get(value, 'date', 'day') ?? 1), 0, 0, 0, 'UTC') ELSE <std::datetime>{}),
      precision := <date::DatePrecision>(value)['precision']
  ), message := ('Failed to parse date with precision from JSON: ' ++ std::to_str(value))));
  CREATE FUNCTION date::rewrite_date(value: tuple<date: std::datetime, precision: date::DatePrecision>) -> OPTIONAL std::datetime USING ((std::datetime_truncate(value.date, 'years') IF (value.precision = date::DatePrecision.Year) ELSE (std::datetime_truncate(value.date, 'months') IF (value.precision = date::DatePrecision.Month) ELSE (std::datetime_truncate(value.date, 'days') IF (value.precision = date::DatePrecision.Day) ELSE <std::datetime>{}))));
  ALTER TYPE events::Event {
      ALTER PROPERTY performed_on {
          CREATE REWRITE
              INSERT 
              USING ((
                  date := date::rewrite_date(__subject__.performed_on),
                  precision := __subject__.performed_on.precision
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  date := date::rewrite_date(__subject__.performed_on),
                  precision := __subject__.performed_on.precision
              ));
      };
  };
  ALTER TYPE occurrence::Identification {
      ALTER PROPERTY identified_on {
          CREATE REWRITE
              INSERT 
              USING ((
                  date := date::rewrite_date(.identified_on),
                  precision := .identified_on.precision
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  date := date::rewrite_date(.identified_on),
                  precision := .identified_on.precision
              ));
      };
  };
  ALTER TYPE samples::Slide {
      ALTER PROPERTY created_on {
          CREATE REWRITE
              INSERT 
              USING ((
                  date := date::rewrite_date(.created_on),
                  precision := .created_on.precision
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  date := date::rewrite_date(.created_on),
                  precision := .created_on.precision
              ));
      };
  };
  ALTER TYPE seq::DNA {
      ALTER PROPERTY extracted_on {
          CREATE REWRITE
              INSERT 
              USING ((
                  date := date::rewrite_date(.extracted_on),
                  precision := .extracted_on.precision
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  date := date::rewrite_date(.extracted_on),
                  precision := .extracted_on.precision
              ));
      };
  };
  ALTER TYPE seq::PCR {
      ALTER PROPERTY performed_on {
          CREATE REWRITE
              INSERT 
              USING ((
                  date := date::rewrite_date(__subject__.performed_on),
                  precision := __subject__.performed_on.precision
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  date := date::rewrite_date(__subject__.performed_on),
                  precision := __subject__.performed_on.precision
              ));
      };
  };
  CREATE FUNCTION default::null_if_empty(s: std::str) -> OPTIONAL std::str USING (WITH
      trimmed := 
          std::str_trim(s)
  SELECT
      (<std::str>{} IF (std::len(trimmed) = 0) ELSE trimmed)
  );
  CREATE TYPE people::User {
      CREATE REQUIRED PROPERTY login: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(16);
          CREATE CONSTRAINT std::min_len_value(5);
      };
      CREATE REQUIRED PROPERTY role: people::UserRole {
          SET default := (people::UserRole.Visitor);
      };
      CREATE REQUIRED PROPERTY email: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY password: std::str {
          CREATE ANNOTATION std::description := 'Password hashing is done within the database, raw password must be used when creating/updating.';
          CREATE REWRITE
              INSERT 
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('bf', 10)) ELSE .password));
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('bf', 10)) ELSE .password));
      };
      CREATE INDEX ON ((.email, .login));
  };
  CREATE GLOBAL default::current_user := (SELECT
      people::User
  FILTER
      (.id = GLOBAL default::current_user_id)
  );
  CREATE TYPE default::Meta {
      CREATE LINK created_by_user: people::User {
          SET default := (SELECT
              GLOBAL default::current_user
          );
      };
      CREATE REQUIRED PROPERTY created: std::datetime {
          SET default := (std::datetime_of_statement());
      };
      CREATE LINK modified_by_user: people::User {
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  GLOBAL default::current_user
              );
      };
      CREATE PROPERTY modified: std::datetime;
      CREATE ANNOTATION std::title := 'Tracking data modifications';
      CREATE PROPERTY lastUpdated := ((.modified ?? .created));
  };
  ALTER TYPE default::Auditable {
      CREATE REQUIRED LINK meta: default::Meta {
          SET default := (INSERT
              default::Meta
          );
          ON SOURCE DELETE DELETE TARGET;
          ON TARGET DELETE RESTRICT;
          CREATE CONSTRAINT std::exclusive;
          CREATE REWRITE
              UPDATE 
              USING (UPDATE
                  .meta
              SET {
                  modified := std::datetime_of_statement()
              });
      };
  };
  CREATE TYPE people::Organisation EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(12);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      CREATE PROPERTY description: std::str {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.description));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.description));
      };
      CREATE REQUIRED PROPERTY kind: people::OrgKind;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(128);
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
  CREATE FUNCTION people::insert_organisation(data: std::json) ->  people::Organisation USING (INSERT
      people::Organisation
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          kind := <people::OrgKind>(data)['kind'],
          description := <std::str>std::json_get(data, 'description')
      });
  CREATE FUNCTION people::insert_or_create_organisation(data: std::json) ->  people::Organisation USING ((IF (std::json_typeof(data) = 'object') THEN (SELECT
      people::insert_organisation(data)
  ) ELSE (IF (std::json_typeof(data) = 'string') THEN (SELECT
      std::assert_exists(people::Organisation FILTER
          (.code = <std::str>data)
      , message := ('Failed to find organisation with code: ' ++ <std::str>data))
  ) ELSE std::assert_exists(<people::Organisation>{}, message := ('Invalid organisation JSON type: ' ++ std::json_typeof(data))))));
  CREATE TYPE people::Person EXTENDING default::Auditable {
      CREATE MULTI LINK organisations: people::Organisation;
      CREATE REQUIRED PROPERTY alias: std::str {
          SET default := (<std::str>{});
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(32);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE PROPERTY comment: std::str;
      CREATE PROPERTY contact: std::str {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.contact));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.contact));
      };
      CREATE REQUIRED PROPERTY first_name: std::str {
          CREATE CONSTRAINT std::max_len_value(30);
          CREATE CONSTRAINT std::min_len_value(1);
      };
      CREATE REQUIRED PROPERTY last_name: std::str {
          CREATE CONSTRAINT std::max_len_value(30);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      ALTER PROPERTY alias {
          CREATE REWRITE
              INSERT 
              USING ((default::null_if_empty(.alias) ?? (WITH
                  default_alias := 
                      std::str_lower(((.first_name)[0] ++ .last_name))
                  ,
                  conflicts := 
                      (SELECT
                          std::count(DETACHED people::Person FILTER
                              (std::str_trim(.alias, '0123456789') = default_alias)
                          )
                      )
                  ,
                  suffix := 
                      (IF (conflicts > 0) THEN <std::str>conflicts ELSE '')
              SELECT
                  (default_alias ++ suffix)
              )));
          CREATE REWRITE
              UPDATE 
              USING ((default::null_if_empty(.alias) ?? (WITH
                  default_alias := 
                      std::str_lower(((.first_name)[0] ++ .last_name))
                  ,
                  conflicts := 
                      (SELECT
                          std::count(DETACHED people::Person FILTER
                              (std::str_trim(.alias, '0123456789') = default_alias)
                          )
                      )
                  ,
                  suffix := 
                      (IF (conflicts > 0) THEN <std::str>conflicts ELSE '')
              SELECT
                  (default_alias ++ suffix)
              )));
      };
      CREATE REQUIRED PROPERTY full_name := (std::array_join([.first_name, .last_name], ' '));
      CREATE INDEX ON ((.alias, .first_name, .last_name));
  };
  CREATE FUNCTION people::insert_person(data: std::json) ->  people::Person USING (INSERT
      people::Person
      {
          first_name := <std::str>(data)['first_name'],
          last_name := <std::str>(data)['last_name'],
          alias := <std::str>std::json_get(data, 'alias'),
          contact := <std::str>std::json_get(data, 'contact'),
          comment := <std::str>std::json_get(data, 'comment'),
          organisations := DISTINCT ((FOR inst IN std::json_array_unpack(std::json_get(data, 'organisations'))
          UNION 
              people::insert_or_create_organisation(inst)))
      });
  CREATE FUNCTION occurrence::externalBiomatByCode(code: OPTIONAL std::str) ->  occurrence::ExternalBioMat USING (SELECT
      std::assert_exists(occurrence::ExternalBioMat FILTER
          (.code = code)
      , message := ('Failed to find external biomaterial with code: ' ++ code))
  );
  CREATE FUNCTION people::personByAlias(alias: std::str) ->  people::Person USING (SELECT
      std::assert_exists(people::Person FILTER
          (.alias = alias)
      , message := ('Failed to find person with alias: ' ++ alias))
  );
  CREATE GLOBAL references::alphabet := ('abcdefghijklmnopqrstuvxyz');
  CREATE TYPE references::Article EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY authors: array<std::str>;
      CREATE REQUIRED PROPERTY year: std::int32 {
          CREATE CONSTRAINT std::min_value(1000);
      };
      CREATE REQUIRED PROPERTY code: std::str {
          SET default := (WITH
              gen_code := 
                  references::generate_article_code(.authors, .year)
              ,
              discriminant := 
                  (SELECT
                      (std::count(DETACHED references::Article FILTER
                          (references::generate_article_code(.authors, .year) = gen_code)
                      ) - 1)
                  )
          SELECT
              (IF (discriminant >= 0) THEN (references::generate_article_code(.authors, .year) ++ (GLOBAL references::alphabet)[discriminant]) ELSE references::generate_article_code(.authors, .year))
          );
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE INDEX ON ((.code, .year));
      CREATE PROPERTY comments: std::str;
      CREATE PROPERTY doi: std::str;
      CREATE PROPERTY journal: std::str;
      CREATE PROPERTY title: std::str;
      CREATE PROPERTY verbatim: std::str;
  };
  CREATE TYPE seq::Gene EXTENDING default::Vocabulary, default::Auditable {
      CREATE REQUIRED PROPERTY motu: std::bool {
          SET default := false;
      };
  };
  CREATE FUNCTION seq::geneByCode(code: std::str) ->  seq::Gene USING (SELECT
      std::assert_exists(seq::Gene FILTER
          (.code = code)
      , message := ('Failed to find gene with code: ' ++ code))
  );
  CREATE TYPE seq::SeqDB EXTENDING default::Auditable, default::Vocabulary {
      CREATE PROPERTY link_template: std::str;
  };
  CREATE FUNCTION seq::seqDbByCode(code: std::str) ->  seq::SeqDB USING (SELECT
      std::assert_exists(seq::SeqDB FILTER
          (.code = code)
      , message := ('Failed to find sequence database with code: ' ++ code))
  );
  CREATE FUNCTION taxonomy::taxonByName(name: std::str) ->  taxonomy::Taxon USING (SELECT
      std::assert_exists(taxonomy::Taxon FILTER
          (.name = name)
      , message := ('Failed to find taxon with name: ' ++ name))
  );
  CREATE ABSTRACT TYPE datasets::AbstractDataset EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK maintainers: people::Person;
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::max_len_value(40);
          CREATE CONSTRAINT std::min_len_value(4);
      };
      CREATE REQUIRED PROPERTY slug: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE ABSTRACT TYPE events::Action EXTENDING default::Auditable {
      CREATE REQUIRED LINK event: events::Event {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
  CREATE TYPE datasets::OccurrenceDataset EXTENDING datasets::AbstractDataset {
      CREATE MULTI LINK occurrences: occurrence::Occurrence {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
      CREATE REQUIRED PROPERTY is_congruent := (std::all(((SELECT
          .occurrences {
              is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true))
          }
      )).is_congruent));
  };
  CREATE TYPE datasets::SiteDataset EXTENDING datasets::AbstractDataset {
      CREATE MULTI LINK sites: location::Site {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
  };
  CREATE TYPE datasets::SeqDataset EXTENDING datasets::AbstractDataset {
      CREATE MULTI LINK sequences: seq::Sequence {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
  };
  CREATE TYPE datasets::DelimitationMethod EXTENDING default::Vocabulary, default::Auditable;
  CREATE TYPE datasets::MOTU EXTENDING default::Auditable {
      CREATE REQUIRED LINK method: datasets::DelimitationMethod;
      CREATE REQUIRED MULTI LINK sequences: seq::Sequence;
      CREATE REQUIRED PROPERTY number: std::int32;
  };
  CREATE TYPE datasets::MOTUDataset EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK generated_by: people::Person;
      CREATE LINK published_in: references::Article;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY label: std::str;
  };
  CREATE ABSTRACT TYPE storage::Storage EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(4);
      };
  };
  CREATE TYPE storage::BioMatStorage EXTENDING storage::Storage;
  CREATE TYPE storage::DNAStorage EXTENDING storage::Storage;
  CREATE ABSTRACT TYPE seq::Primer EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY sequence: std::str {
          CREATE CONSTRAINT std::min_len_value(5);
      };
  };
  CREATE TYPE sampling::Habitat EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE INDEX ON (.label);
      CREATE PROPERTY description: std::str;
  };
  CREATE TYPE seq::PCRReversePrimer EXTENDING seq::Primer, default::Auditable;
  CREATE TYPE samples::Fixative EXTENDING default::Vocabulary, default::Auditable {
      CREATE ANNOTATION std::description := 'Describes a conservation method for a sample.';
  };
  CREATE TYPE events::SamplingMethod EXTENDING default::Vocabulary, default::Auditable;
  CREATE TYPE events::AbioticParameter EXTENDING default::Vocabulary, default::Auditable {
      CREATE REQUIRED PROPERTY unit: std::str;
  };
  CREATE TYPE events::Program EXTENDING default::Auditable {
      CREATE PROPERTY end_year: std::int32;
      CREATE PROPERTY start_year: std::int32 {
          CREATE CONSTRAINT std::min_value(1900);
      };
      CREATE CONSTRAINT std::expression ON ((.start_year <= .end_year));
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE INDEX ON ((.code, .label));
      CREATE MULTI LINK funding_agencies: people::Organisation;
      CREATE REQUIRED MULTI LINK managers: people::Person;
      CREATE PROPERTY description: std::str;
  };
  CREATE TYPE storage::SlideStorage EXTENDING storage::Storage;
  CREATE TYPE seq::ChromatoPrimer EXTENDING seq::Primer, default::Auditable;
  CREATE TYPE seq::SequencingInstitute EXTENDING default::Auditable {
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE seq::Chromatogram EXTENDING default::Auditable {
      CREATE REQUIRED LINK primer: seq::ChromatoPrimer;
      CREATE REQUIRED PROPERTY YAS_number: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE REWRITE
              INSERT 
              USING ((.code ?? (SELECT
                  ((.YAS_number ++ '|') ++ .primer.code)
              )));
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED LINK PCR: seq::PCR;
      CREATE REQUIRED LINK provider: seq::SequencingInstitute;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY quality: seq::ChromatoQuality;
  };
  CREATE TYPE seq::BatchRequest EXTENDING default::Auditable {
      CREATE REQUIRED LINK requested_by: people::Person;
      CREATE MULTI LINK requested_to: people::Person;
      CREATE REQUIRED MULTI LINK content: seq::DNA;
      CREATE REQUIRED LINK target_gene: seq::Gene;
      CREATE PROPERTY achieved_on: std::datetime;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY label: std::str;
      CREATE PROPERTY requested_on: std::datetime {
          CREATE ANNOTATION std::description := 'If empty, the request is a draft and can not be processed yet.';
      };
  };
  CREATE TYPE samples::ContentType EXTENDING default::Vocabulary, default::Auditable;
  CREATE TYPE sampling::HabitatGroup EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY exclusive_elements: std::bool {
          SET default := true;
      };
      CREATE LINK depends: sampling::Habitat;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE INDEX ON (.label);
  };
  CREATE TYPE events::AbioticMeasurement EXTENDING events::Action {
      CREATE REQUIRED LINK param: events::AbioticParameter {
          ON TARGET DELETE DELETE SOURCE;
      };
      CREATE CONSTRAINT std::exclusive ON ((.event, .param));
      CREATE REQUIRED PROPERTY value: std::float32;
  };
  ALTER TYPE samples::Sample {
      CREATE REQUIRED LINK conservation: samples::Fixative;
      CREATE REQUIRED LINK type: samples::ContentType;
      CREATE PROPERTY tube := ((.type.code ++ <std::str>.number));
      ALTER PROPERTY number {
          CREATE REWRITE
              INSERT 
              USING (SELECT
                  (1 + (std::max(samples::Sample.number FILTER
                      ((samples::Sample.biomat = __subject__.biomat) AND (samples::Sample.type = __subject__.type))
                  ) ?? 0))
              );
      };
  };
  CREATE TYPE samples::BundledSpecimens EXTENDING samples::Sample {
      CREATE ANNOTATION std::description := 'A tube containing several specimens.';
      CREATE REQUIRED PROPERTY quantity: std::int16 {
          CREATE CONSTRAINT std::min_value(2);
      };
  };
  CREATE TYPE seq::PCRForwardPrimer EXTENDING seq::Primer, default::Auditable;
  CREATE TYPE seq::DNAExtractionMethod EXTENDING default::Vocabulary, default::Auditable;
  CREATE TYPE events::Sampling EXTENDING events::Action {
      CREATE REQUIRED PROPERTY number: events::SamplingNumber {
          SET readonly := true;
      };
      CREATE REQUIRED SINGLE PROPERTY code := (WITH
          id := 
              .id
          ,
          event := 
              .event
          ,
          siblings := 
              (SELECT
                  DETACHED events::Sampling
              FILTER
                  (.event.code = event.code)
              ORDER BY
                  .number ASC
              )
          ,
          rank := 
              ((SELECT
                  (std::assert_single(std::enumerate(siblings) FILTER
                      (.1.id = id)
                  )).0
              ) ?? std::count(siblings))
          ,
          suffix := 
              (IF (std::count(siblings) > 1) THEN ('.' ++ <std::str>(rank + 1)) ELSE '')
      SELECT
          std::assert_single(std::assert_exists((event.code ++ suffix), message := (((("Failed to generate sampling code. Event code: '" ++ event.code) ++ "'. Suffix: '") ++ suffix) ++ "'")))
      );
      CREATE MULTI LINK fixatives: samples::Fixative;
      CREATE MULTI LINK habitats: sampling::Habitat;
      CREATE MULTI LINK methods: events::SamplingMethod;
      CREATE MULTI LINK target_taxa: taxonomy::Taxon {
          CREATE ANNOTATION std::title := 'Taxonomic groups that were the target of the sampling effort';
      };
      CREATE MULTI PROPERTY access_points: std::str;
      CREATE PROPERTY comments: std::str;
      CREATE PROPERTY sampling_duration: std::int32;
      CREATE REQUIRED PROPERTY sampling_target: events::SamplingTarget;
  };
  CREATE ABSTRACT TYPE traits::AbstractTrait {
      CREATE REQUIRED PROPERTY category: traits::Category;
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY name: std::str;
      CREATE REQUIRED MULTI PROPERTY scopes: traits::TraitDefinitionScope {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE CONSTRAINT std::exclusive ON ((.category, .name));
  };
  CREATE TYPE traits::QualitativeTrait EXTENDING traits::AbstractTrait, default::Auditable {
      CREATE REQUIRED PROPERTY value: std::str;
      CREATE CONSTRAINT std::exclusive ON ((.name, .value));
  };
  CREATE TYPE seq::PCRSpecificity EXTENDING default::Auditable {
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE events::Event {
      CREATE MULTI LINK abiotic_measurements := (.<event[IS events::AbioticMeasurement]);
      CREATE REQUIRED MULTI LINK performed_by: people::Person;
      CREATE MULTI LINK programs: events::Program;
      CREATE MULTI LINK samplings := (.<event[IS events::Sampling]);
  };
  ALTER TYPE location::Site {
      CREATE MULTI LINK datasets := (.<sites[IS datasets::SiteDataset]);
      CREATE MULTI LINK events := (.<site[IS events::Event]);
  };
  ALTER TYPE occurrence::Occurrence {
      CREATE REQUIRED LINK sampling: events::Sampling;
      CREATE MULTI LINK published_in: references::Article {
          CREATE PROPERTY original_source: std::bool;
      };
  };
  ALTER TYPE occurrence::BioMaterial {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.code THEN .code ELSE (((.identification.taxon.code ++ '[') ++ .sampling.code) ++ ']')));
      };
  };
  ALTER TYPE occurrence::Identification {
      CREATE LINK identified_by: people::Person;
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE MULTI LINK bundles := (.<biomat[IS samples::BundledSpecimens]);
  };
  ALTER TYPE samples::Slide {
      CREATE REQUIRED MULTI LINK mounted_by: people::Person;
      CREATE REQUIRED LINK storage: storage::SlideStorage;
      CREATE CONSTRAINT std::exclusive ON ((.storage, .storage_position));
  };
  ALTER TYPE samples::Specimen {
      ALTER PROPERTY molecular_number {
          CREATE REWRITE
              INSERT 
              USING (<std::str>(SELECT
                  std::count(samples::Specimen FILTER
                      (samples::Specimen.biomat.sampling.event.site = __subject__.biomat.sampling.event.site)
                  )
              ));
      };
      ALTER PROPERTY morphological_code {
          CREATE REWRITE
              INSERT 
              USING ((((__subject__.biomat.code ++ '[') ++ .tube) ++ ']'));
          CREATE REWRITE
              UPDATE 
              USING ((((__subject__.biomat.code ++ '[') ++ .tube) ++ ']'));
      };
      CREATE LINK dissected_by: people::Person;
      CREATE MULTI LINK slides := (.<specimen[IS samples::Slide]);
  };
  ALTER TYPE seq::Sequence {
      CREATE REQUIRED LINK sampling: events::Sampling;
      CREATE REQUIRED PROPERTY category := (std::assert_exists((IF (__source__ IS seq::AssembledSequence) THEN occurrence::OccurrenceCategory.Internal ELSE (IF (__source__ IS seq::ExternalSequence) THEN occurrence::OccurrenceCategory.External ELSE {})), message := (('Occurrence category for seq::Sequence subtype ' ++ __source__.__type__.name) ++ ' is undefined')));
      CREATE REQUIRED LINK gene: seq::Gene;
  };
  ALTER TYPE seq::AssembledSequence {
      ALTER LINK sampling {
          SET OWNED;
          SET REQUIRED;
          SET TYPE events::Sampling;
          CREATE REWRITE
              INSERT 
              USING (SELECT
                  .specimen.biomat.sampling
              );
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  .specimen.biomat.sampling
              );
      };
      CREATE REQUIRED MULTI LINK assembled_by: people::Person;
      CREATE MULTI LINK published_in: references::Article;
      CREATE REQUIRED MULTI LINK chromatograms: seq::Chromatogram;
  };
  ALTER TYPE seq::DNA {
      CREATE REQUIRED MULTI LINK extracted_by: people::Person;
      CREATE MULTI LINK PCRs := (.<DNA[IS seq::PCR]);
      CREATE REQUIRED LINK extraction_method: seq::DNAExtractionMethod;
      CREATE REQUIRED LINK stored_in: storage::DNAStorage;
  };
  ALTER TYPE seq::PCR {
      CREATE MULTI LINK chromatograms := (.<PCR[IS seq::Chromatogram]);
      CREATE REQUIRED LINK gene: seq::Gene;
      CREATE REQUIRED LINK forward_primer: seq::PCRForwardPrimer;
      CREATE REQUIRED LINK reverse_primer: seq::PCRReversePrimer;
  };
  CREATE ABSTRACT TYPE tokens::Token {
      CREATE REQUIRED PROPERTY expires: std::datetime;
      CREATE REQUIRED PROPERTY token: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE tokens::UserInvitation EXTENDING tokens::Token {
      CREATE REQUIRED LINK issued_by: people::User {
          SET default := (GLOBAL default::current_user);
      };
      CREATE REQUIRED LINK identity: people::Person;
      CREATE REQUIRED PROPERTY email: std::str;
      CREATE INDEX ON (.email);
      CREATE REQUIRED PROPERTY role: people::UserRole;
  };
  CREATE TYPE admin::EmailSettings {
      CREATE REQUIRED PROPERTY from_address: std::str;
      CREATE REQUIRED PROPERTY from_name: std::str;
      CREATE REQUIRED PROPERTY host: std::str;
      CREATE REQUIRED PROPERTY password: std::str;
      CREATE REQUIRED PROPERTY port: std::int32;
      CREATE REQUIRED PROPERTY user: std::str;
  };
  CREATE TYPE admin::InstanceSettings {
      CREATE REQUIRED PROPERTY allow_contributor_signup: std::bool {
          SET default := true;
          CREATE ANNOTATION std::description := 'Whether a user can request a new account with contributor privileges';
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::max_len_value(20);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE REQUIRED PROPERTY public: std::bool {
          SET default := true;
          CREATE ANNOTATION std::description := 'Whether parts of the platform are open to the public (anonymous users).';
      };
  };
  CREATE TYPE admin::SecuritySettings {
      CREATE REQUIRED PROPERTY invitation_token_lifetime: std::int32 {
          SET default := 7;
          CREATE ANNOTATION std::description := 'Validity period for an account invitation in days';
          CREATE CONSTRAINT std::min_value(1);
      };
      CREATE REQUIRED PROPERTY jwt_secret_key: std::str {
          CREATE CONSTRAINT std::min_len_value(32);
      };
      CREATE REQUIRED PROPERTY min_password_strength: std::int32 {
          SET default := 3;
          CREATE CONSTRAINT std::max_value(5);
          CREATE CONSTRAINT std::min_value(3);
      };
      CREATE REQUIRED PROPERTY refresh_token_lifetime: std::int32 {
          SET default := ((24 * 30));
          CREATE ANNOTATION std::description := 'Validity period for a session refresh token in hours';
          CREATE CONSTRAINT std::min_value(1);
      };
  };
  CREATE TYPE admin::Settings {
      CREATE LINK email := (SELECT
          admin::EmailSettings 
      LIMIT
          1
      );
      CREATE REQUIRED LINK instance := (std::assert_exists((SELECT
          admin::InstanceSettings 
      LIMIT
          1
      ), message := 'Instance settings are not intialized. This is a fatal error and should never occur if database was properly initialized.'));
      CREATE REQUIRED LINK security := (std::assert_exists((SELECT
          admin::SecuritySettings 
      LIMIT
          1
      ), message := 'Security settings are not intialized. This is a fatal error and should never occur if database was properly initialized.'));
      CREATE REQUIRED LINK superadmin: people::User;
      CREATE PROPERTY geoapify_api_key: std::str;
  };
  CREATE TYPE admin::GeoapifyUsage {
      CREATE REQUIRED PROPERTY date: std::cal::local_date {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY requests: std::int32;
  };
  ALTER TYPE datasets::MOTU {
      CREATE REQUIRED LINK dataset: datasets::MOTUDataset;
  };
  ALTER TYPE datasets::MOTUDataset {
      CREATE MULTI LINK MOTUs := (.<dataset[IS datasets::MOTU]);
  };
  ALTER TYPE datasets::OccurrenceDataset {
      CREATE MULTI LINK sites := (.occurrences.sampling.event.site);
  };
  ALTER TYPE datasets::SeqDataset {
      CREATE MULTI LINK sites := (.sequences.sampling.event.site);
  };
  ALTER TYPE people::User {
      CREATE REQUIRED LINK identity: people::Person {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE RESTRICT;
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE default::Meta {
      CREATE PROPERTY created_by := ((
          id := .created_by_user.id,
          login := .created_by_user.login,
          name := .created_by_user.identity.full_name,
          alias := .created_by_user.identity.alias
      ));
      CREATE PROPERTY updated_by := ((
          id := .modified_by_user.id,
          login := .modified_by_user.login,
          name := .modified_by_user.identity.full_name,
          alias := .modified_by_user.identity.alias
      ));
  };
  CREATE TYPE seq::SeqReference {
      CREATE REQUIRED LINK db: seq::SeqDB {
          ON TARGET DELETE DELETE SOURCE;
      };
      CREATE REQUIRED PROPERTY accession: std::str {
          CREATE CONSTRAINT std::max_len_value(32);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE REQUIRED PROPERTY code := (((.db.code ++ ':') ++ .accession));
      CREATE REQUIRED PROPERTY is_origin: std::bool {
          SET default := false;
      };
      CREATE CONSTRAINT std::exclusive ON ((.db, .accession));
  };
  ALTER TYPE seq::ExternalSequence {
      CREATE MULTI LINK referenced_in: seq::SeqReference {
          CREATE CONSTRAINT std::exclusive;
      };
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (WITH
                  suffix := 
                      (IF (.origin = seq::ExtSeqOrigin.Lab) THEN 'lab' ELSE (IF (.origin = seq::ExtSeqOrigin.PersCom) THEN 'perscom' ELSE (WITH
                          sources := 
                              ((SELECT
                                  .referenced_in
                              FILTER
                                  .is_origin
                              )).code
                      SELECT
                          std::array_join(std::array_agg(sources), '|')
                      )))
              SELECT
                  ((((((.identification.taxon.code ++ '[') ++ .sampling.code) ++ ']') ++ .specimen_identifier) ++ '|') ++ suffix)
              );
      };
  };
  ALTER TYPE events::Sampling {
      CREATE MULTI LINK external_seqs := (.<sampling[IS seq::ExternalSequence]);
      CREATE MULTI LINK samples := (.<sampling[IS occurrence::BioMaterial]);
      CREATE MULTI LINK occurring_taxa := (WITH
          ext_samples_no_seqs := 
              (SELECT
                  .samples[IS occurrence::ExternalBioMat]
              FILTER
                  NOT (EXISTS ([IS occurrence::ExternalBioMat].sequences))
              )
      SELECT
          DISTINCT (((ext_samples_no_seqs.identification.taxon UNION .external_seqs.identification.taxon) UNION .samples[IS occurrence::InternalBioMat].identified_taxa))
      );
  };
  ALTER TYPE people::Organisation {
      CREATE MULTI LINK people := (.<organisations[IS people::Person]);
  };
  CREATE TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY email: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE INDEX ON (.email);
      CREATE REQUIRED PROPERTY created_on: std::datetime {
          CREATE REWRITE
              INSERT 
              USING (std::datetime_of_statement());
      };
      CREATE REQUIRED PROPERTY email_verified: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY first_name: std::str;
      CREATE REQUIRED PROPERTY last_name: std::str;
      CREATE REQUIRED PROPERTY full_name := (((.first_name ++ ' ') ++ .last_name));
      CREATE PROPERTY motive: std::str;
      CREATE PROPERTY organisation: std::str;
  };
  CREATE TYPE tokens::EmailVerification EXTENDING tokens::Token {
      CREATE REQUIRED LINK user_request: people::PendingUserRequest {
          ON TARGET DELETE DELETE SOURCE;
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE people::Person {
      CREATE LINK user := (.<identity[IS people::User]);
      CREATE OPTIONAL PROPERTY role := (.user.role);
  };
  CREATE TYPE storage::Collection {
      CREATE REQUIRED LINK maintainers: people::Person;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(4);
      };
      CREATE REQUIRED LINK taxon: taxonomy::Taxon;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(4);
      };
  };
  CREATE TYPE tokens::PasswordReset EXTENDING tokens::Token {
      CREATE REQUIRED LINK user: people::User {
          ON TARGET DELETE DELETE SOURCE;
          CREATE DELEGATED CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE tokens::SessionRefreshToken EXTENDING tokens::Token {
      CREATE REQUIRED LINK user: people::User {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
  ALTER TYPE storage::Storage {
      CREATE REQUIRED LINK collection: storage::Collection;
  };
  ALTER TYPE samples::Slide {
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (SELECT
                  std::array_join([.storage.collection.code, .storage.code, <std::str>.storage_position], '_')
              );
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  std::array_join([.storage.collection.code, .storage.code, <std::str>.storage_position], '_')
              );
      };
  };
  ALTER TYPE sampling::Habitat {
      CREATE REQUIRED LINK in_group: sampling::HabitatGroup {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
  ALTER TYPE sampling::HabitatGroup {
      CREATE LINK elements := (.<in_group[IS sampling::Habitat]);
  };
  ALTER TYPE sampling::Habitat {
      CREATE LINK incompatible := ((.in_group.elements IF .in_group.exclusive_elements ELSE {}));
  };
  ALTER TYPE seq::Chromatogram {
      CREATE MULTI LINK sequences := (.<chromatograms[IS seq::AssembledSequence]);
  };
  ALTER TYPE storage::Collection {
      CREATE MULTI LINK bio_mat_storages := (.<collection[IS storage::BioMatStorage]);
  };
  ALTER TYPE storage::Collection {
      CREATE MULTI LINK DNA_storages := (.<collection[IS storage::DNAStorage]);
      CREATE MULTI LINK slide_storages := (.<collection[IS storage::SlideStorage]);
  };
};
