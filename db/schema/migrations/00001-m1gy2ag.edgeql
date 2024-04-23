CREATE MIGRATION m1gy2agbdoat3bsdqbe7apdaaw6n7poybyu7fo72zta4lxy4imigvq
    ONTO initial
{
  CREATE EXTENSION pgcrypto VERSION '1.3';
  CREATE EXTENSION auth VERSION '1.0';
  CREATE MODULE admin IF NOT EXISTS;
  CREATE MODULE datasets IF NOT EXISTS;
  CREATE MODULE date IF NOT EXISTS;
  CREATE MODULE event IF NOT EXISTS;
  CREATE MODULE location IF NOT EXISTS;
  CREATE MODULE occurrence IF NOT EXISTS;
  CREATE MODULE people IF NOT EXISTS;
  CREATE MODULE reference IF NOT EXISTS;
  CREATE MODULE samples IF NOT EXISTS;
  CREATE MODULE seq IF NOT EXISTS;
  CREATE MODULE storage IF NOT EXISTS;
  CREATE MODULE taxonomy IF NOT EXISTS;
  CREATE GLOBAL default::current_user_id -> std::uuid {
      SET default := (<std::uuid>{});
  };
  CREATE SCALAR TYPE people::UserRole EXTENDING enum<Visitor, Contributor, Maintainer, Admin>;
  CREATE ABSTRACT TYPE default::Auditable {
      CREATE ANNOTATION std::title := 'Auto-generation of timestamps';
  };
  CREATE TYPE people::Person EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY alias: std::str {
          SET default := (<std::str>{});
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(32);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE REQUIRED PROPERTY first_name: std::str {
          CREATE CONSTRAINT std::max_len_value(30);
          CREATE CONSTRAINT std::min_len_value(1);
      };
      CREATE REQUIRED PROPERTY last_name: std::str {
          CREATE CONSTRAINT std::max_len_value(30);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      CREATE PROPERTY contact: std::str;
      CREATE REQUIRED PROPERTY full_name := (std::array_join(std::array_agg({.first_name, .last_name}), ' '));
      CREATE PROPERTY comment: std::str;
  };
  CREATE TYPE people::User {
      CREATE LINK identity: people::Person {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE RESTRICT;
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY email_confirmed: std::bool {
          SET default := false;
      };
      CREATE PROPERTY is_active := ((.email_confirmed AND EXISTS (.identity)));
      CREATE REQUIRED PROPERTY email: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY login: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(16);
          CREATE CONSTRAINT std::min_len_value(5);
      };
      CREATE REQUIRED PROPERTY password: std::str {
          CREATE ANNOTATION std::description := 'Password hashing is done within the database, raw password must be used when creating/updating.';
          CREATE REWRITE
              INSERT 
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des')) ELSE .password));
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des')) ELSE .password));
      };
      CREATE REQUIRED PROPERTY role: people::UserRole {
          SET default := (people::UserRole.Visitor);
      };
  };
  CREATE ALIAS people::ActiveUser := (
      SELECT
          people::User
      FILTER
          .is_active
  );
  CREATE ALIAS people::InactiveUser := (
      SELECT
          people::User
      FILTER
          NOT (.is_active)
  );
  CREATE ABSTRACT ANNOTATION default::example;
  CREATE GLOBAL default::current_user := (SELECT
      people::User
  FILTER
      (.id = GLOBAL default::current_user_id)
  );
  CREATE SCALAR TYPE date::DatePrecision EXTENDING enum<Year, Month, Day, Unknown>;
  CREATE ABSTRACT CONSTRAINT date::required_unless_unknown(date: std::datetime, precision: date::DatePrecision) {
      SET errmessage := "Date value is required, except when precision is 'Unknown'";
      USING ((EXISTS (date) IF (precision != date::DatePrecision.Unknown) ELSE NOT (EXISTS (date))));
  };
  CREATE SCALAR TYPE event::SamplingTarget EXTENDING enum<Community, Unknown, Taxa>;
  CREATE ABSTRACT TYPE event::Event EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK performed_by: people::Person;
      CREATE REQUIRED PROPERTY performed_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(__subject__.date, __subject__.precision);
      };
  };
  CREATE TYPE event::Sampling EXTENDING event::Event {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE ANNOTATION default::example := 'SOMESITE_202301 ; SOMESITE_202301_1';
          CREATE ANNOTATION std::description := 'Format : SITE_YEARMONTH_NUMBER. The NUMBER suffix is not appended if the site and month tuple is unique.';
          CREATE ANNOTATION std::title := 'Unique sampling identifier, auto-generated at sampling creation.';
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY is_donation: std::bool;
      CREATE PROPERTY sampling_duration: std::duration;
      CREATE REQUIRED PROPERTY sampling_target: event::SamplingTarget;
  };
  CREATE TYPE seq::Gene EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE ANNOTATION default::example := 'COI';
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE ANNOTATION default::example := 'cytochrome oxydase';
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY motu: std::bool;
  };
  CREATE TYPE event::AbioticMeasurement EXTENDING event::Event {
      CREATE LINK related_sampling: event::Sampling;
      CREATE REQUIRED PROPERTY value: std::float32;
  };
  CREATE TYPE event::Spotting EXTENDING event::Event;
  CREATE TYPE occurrence::Identification EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY identified_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(__subject__.date, __subject__.precision);
      };
      CREATE REQUIRED LINK identified_by: people::Person;
  };
  CREATE ABSTRACT TYPE occurrence::Occurrence EXTENDING default::Auditable {
      CREATE REQUIRED LINK identification: occurrence::Identification {
          ON SOURCE DELETE DELETE TARGET;
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED LINK sampling: event::Sampling;
      CREATE PROPERTY comments: std::str;
  };
  CREATE TYPE samples::BioMaterial EXTENDING occurrence::Occurrence {
      CREATE REQUIRED PROPERTY created_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE ANNOTATION std::description := "Format like 'taxon_code|sampling_code'";
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED MULTI LINK sorted_by: people::Person;
  };
  CREATE TYPE default::Picture {
      CREATE PROPERTY legend: std::str;
      CREATE REQUIRED PROPERTY path: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE samples::Slide EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY created_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE MULTI LINK pictures: default::Picture {
          ON SOURCE DELETE DELETE TARGET;
      };
      CREATE REQUIRED MULTI LINK mounted_by: people::Person;
      CREATE REQUIRED PROPERTY storage_position: std::int16;
      CREATE PROPERTY code: std::str {
          CREATE ANNOTATION std::description := "Generated as '{collectionCode}_{containerCode}_{slidePositionInBox}'";
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comment: std::str;
      CREATE REQUIRED PROPERTY label: std::str;
  };
  CREATE SCALAR TYPE seq::DNAQuality EXTENDING enum<Unknown, Contaminated, Bad, Good>;
  CREATE TYPE seq::DNA EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY extracted_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE REQUIRED MULTI LINK extracted_by: people::Person;
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
  CREATE SCALAR TYPE seq::PCRQuality EXTENDING enum<Failure, Acceptable, Good, Unknown>;
  CREATE TYPE seq::PCR EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY performed_on: tuple<date: std::datetime, precision: date::DatePrecision> {
          CREATE CONSTRAINT date::required_unless_unknown(.date, .precision);
      };
      CREATE REQUIRED LINK DNA: seq::DNA;
      CREATE REQUIRED LINK gene: seq::Gene;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY quality: seq::PCRQuality;
  };
  CREATE FUNCTION date::rewrite_date(value: tuple<date: std::datetime, precision: date::DatePrecision>) -> OPTIONAL std::datetime USING ((std::datetime_truncate(value.date, 'years') IF (value.precision = date::DatePrecision.Year) ELSE (std::datetime_truncate(value.date, 'months') IF (value.precision = date::DatePrecision.Month) ELSE (std::datetime_truncate(value.date, 'days') IF (value.precision = date::DatePrecision.Day) ELSE <std::datetime>{}))));
  ALTER TYPE event::Event {
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
  ALTER TYPE samples::BioMaterial {
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
  ALTER TYPE people::Person {
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
      ALTER PROPERTY contact {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.contact));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.contact));
      };
  };
  CREATE SCALAR TYPE people::InstitutionKind EXTENDING enum<Lab, FundingAgency, SequencingPlatform, Other>;
  CREATE TYPE people::Institution EXTENDING default::Auditable {
      CREATE PROPERTY description: std::str {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.description));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.description));
      };
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(12);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      CREATE INDEX ON (.code);
      CREATE REQUIRED PROPERTY kind: people::InstitutionKind;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(128);
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
  CREATE ABSTRACT TYPE seq::Sequence {
      CREATE REQUIRED LINK gene: seq::Gene;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comments: std::str;
      CREATE PROPERTY legacy: tuple<id: std::int32, code: std::str, alignment_code: std::str> {
          CREATE ANNOTATION std::description := 'Legacy identifiers for retrocompatibility with data stored in GOTIT.';
      };
  };
  CREATE TYPE datasets::Alignment EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK sequences: seq::Sequence;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY label: std::str;
  };
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
      CREATE PROPERTY created_by := ((
          name := .created_by_user.identity.full_name,
          alias := .created_by_user.identity.alias
      ));
      CREATE PROPERTY updated_by := ((
          name := .modified_by_user.identity.full_name,
          alias := .modified_by_user.identity.alias
      ));
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
  CREATE ABSTRACT TYPE default::Vocabulary {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(8);
          CREATE CONSTRAINT std::min_len_value(2);
          CREATE ANNOTATION std::title := 'An expressive, unique, user-generated uppercase alphanumeric code';
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(4);
      };
      CREATE ANNOTATION std::title := 'An extensible list of terms';
  };
  CREATE TYPE datasets::DelimitationMethod EXTENDING default::Vocabulary, default::Auditable;
  CREATE TYPE datasets::MOTU EXTENDING default::Auditable {
      CREATE REQUIRED LINK method: datasets::DelimitationMethod;
      CREATE REQUIRED MULTI LINK sequences: seq::Sequence;
      CREATE REQUIRED PROPERTY number: std::int16;
  };
  CREATE TYPE reference::Article EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY authors: std::str;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY title: std::str;
      CREATE REQUIRED PROPERTY year: std::int16 {
          CREATE CONSTRAINT std::min_value(1500);
      };
  };
  CREATE TYPE datasets::MOTUDataset EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK generated_by: people::Person;
      CREATE LINK published_in: reference::Article;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY label: std::str;
  };
  CREATE TYPE event::Program EXTENDING default::Auditable {
      CREATE PROPERTY end_year: std::int16;
      CREATE PROPERTY start_year: std::int16 {
          CREATE CONSTRAINT std::min_value(1900);
      };
      CREATE CONSTRAINT std::expression ON ((.start_year <= .end_year));
      CREATE MULTI LINK funding_agencies: people::Institution;
      CREATE REQUIRED MULTI LINK managers: people::Person;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
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
  CREATE TYPE seq::DNAExtractionMethod EXTENDING default::Vocabulary, default::Auditable;
  CREATE SCALAR TYPE taxonomy::Rank EXTENDING enum<Kingdom, Phylum, Class, Order, Family, Genus, Species, Subspecies>;
  CREATE SCALAR TYPE taxonomy::TaxonStatus EXTENDING enum<Accepted, Synonym, Unclassified>;
  CREATE TYPE taxonomy::Taxon EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::min_len_value(4);
      };
      CREATE REQUIRED PROPERTY status: taxonomy::TaxonStatus;
      CREATE CONSTRAINT std::exclusive ON ((.name, .status));
      CREATE REQUIRED PROPERTY rank: taxonomy::Rank;
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) = 3)) EXCEPT ((.rank != taxonomy::Rank.Subspecies));
      CREATE CONSTRAINT std::expression ON (NOT (std::contains(.name, ' '))) EXCEPT (((.rank = taxonomy::Rank.Species) OR (.rank = taxonomy::Rank.Subspecies)));
      CREATE LINK parent: taxonomy::Taxon {
          ON TARGET DELETE DELETE SOURCE;
      };
      CREATE CONSTRAINT std::expression ON (EXISTS (.parent)) EXCEPT ((.rank = taxonomy::Rank.Kingdom));
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) = 2)) EXCEPT ((.rank != taxonomy::Rank.Species));
      CREATE INDEX ON (.name);
      CREATE INDEX ON (.rank);
      CREATE INDEX ON (.status);
      CREATE MULTI LINK children := (.<parent[IS taxonomy::Taxon]);
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
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (WITH
                  chopped := 
                      std::str_split(.name, ' ')
                  ,
                  suffix := 
                      ('[syn]' IF (.status = taxonomy::TaxonStatus.Synonym) ELSE '')
              SELECT
                  (IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE (IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) THEN (.name ++ suffix) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix)))
              );
          CREATE REWRITE
              UPDATE 
              USING (WITH
                  chopped := 
                      std::str_split(.name, ' ')
                  ,
                  suffix := 
                      ('[syn]' IF (.status = taxonomy::TaxonStatus.Synonym) ELSE '')
              SELECT
                  (IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE (IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) THEN (.name ++ suffix) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix)))
              );
      };
      CREATE PROPERTY comment: std::str;
  };
  CREATE TYPE default::Conservation EXTENDING default::Vocabulary, default::Auditable {
      CREATE ANNOTATION std::description := 'Describes a conservation method for a sample.';
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
  CREATE TYPE location::SiteDataset EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK maintainers: people::Person;
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::max_len_value(40);
          CREATE CONSTRAINT std::min_len_value(4);
      };
  };
  CREATE TYPE event::AbioticParameter EXTENDING default::Vocabulary, default::Auditable {
      CREATE REQUIRED PROPERTY unit: std::str;
  };
  CREATE TYPE samples::ContentType EXTENDING default::Vocabulary, default::Auditable;
  CREATE TYPE seq::AssembledSequence EXTENDING seq::Sequence, default::Auditable {
      CREATE REQUIRED LINK identification: occurrence::Identification;
      CREATE REQUIRED PROPERTY is_reference: std::bool {
          CREATE ANNOTATION std::description := 'Whether this sequence should be used as the reference for the identification of the specimen';
      };
      CREATE REQUIRED MULTI LINK assembled_by: people::Person;
      CREATE MULTI LINK published_in: reference::Article;
      CREATE PROPERTY accession_number: std::str {
          CREATE ANNOTATION std::description := 'The NCBI accession number, if the sequence was uploaded.';
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY alignmentCode: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE storage::SlideStorage EXTENDING storage::Storage;
  CREATE TYPE location::HabitatTag EXTENDING default::Auditable {
      CREATE LINK parent: location::HabitatTag;
      CREATE PROPERTY color: std::str {
          CREATE REWRITE
              INSERT 
              USING (((.color ?? .parent.color) IF EXISTS (.parent) ELSE <std::str>{}));
          CREATE REWRITE
              UPDATE 
              USING (((.color ?? .parent.color) IF EXISTS (.parent) ELSE <std::str>{}));
      };
      CREATE PROPERTY description: std::str;
      CREATE PROPERTY is_required: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY label: std::str;
  };
  CREATE SCALAR TYPE location::CoordinateMaxPrecision EXTENDING enum<m10, m100, Km1, Km10, Km100, Unknown>;
  CREATE TYPE location::Site EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(8);
          CREATE CONSTRAINT std::min_len_value(4);
          CREATE ANNOTATION std::description := 'A short, unique, user-generated, alphanumeric identifier. Recommended size is 8.';
          CREATE ANNOTATION std::title := 'Site identifier';
      };
      CREATE REQUIRED MULTI LINK habitat_tags: location::HabitatTag {
          ON TARGET DELETE RESTRICT;
          CREATE ANNOTATION std::title := 'A list of descriptors for the habitat that was targeted.';
      };
      CREATE PROPERTY altitudeRange: tuple<min: std::float32, max: std::float32> {
          CREATE ANNOTATION std::title := 'The site elevation in meters';
      };
      CREATE REQUIRED PROPERTY coordinates: tuple<precision: location::CoordinateMaxPrecision, latitude: std::float32, longitude: std::float32> {
          CREATE CONSTRAINT std::expression ON (((((.latitude <= 90) AND (.latitude >= -90)) AND (.longitude <= 180)) AND (.longitude >= -180)));
      };
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE event::SamplingMethod EXTENDING default::Vocabulary, default::Auditable;
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
  CREATE TYPE location::Country EXTENDING default::Auditable {
      CREATE ANNOTATION std::description := 'Countries as defined in the ISO 3166-1 norm.';
      CREATE REQUIRED PROPERTY code: std::str {
          SET readonly := true;
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(2);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      CREATE REQUIRED PROPERTY name: std::str {
          SET readonly := true;
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE seq::PCRForwardPrimer EXTENDING seq::Primer, default::Auditable;
  CREATE TYPE storage::DNAStorage EXTENDING storage::Storage;
  CREATE ABSTRACT TYPE samples::Sample EXTENDING default::Auditable {
      CREATE REQUIRED LINK biomat: samples::BioMaterial;
      CREATE REQUIRED LINK conservation: default::Conservation;
      CREATE REQUIRED LINK type: samples::ContentType;
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY number: std::int16 {
          CREATE REWRITE
              INSERT 
              USING (SELECT
                  (1 + (std::max(samples::Sample.number FILTER
                      ((samples::Sample.biomat = __subject__.biomat) AND (samples::Sample.type = __subject__.type))
                  ) ?? 0))
              );
          CREATE ANNOTATION std::description := 'Incremental number that discriminates between tubes having the same type in a bio material lot. Used to generate the tube code.';
      };
      CREATE PROPERTY tube := ((.type.code ++ <std::str>.number));
  };
  CREATE TYPE samples::Specimen EXTENDING samples::Sample {
      CREATE MULTI LINK pictures: default::Picture {
          ON SOURCE DELETE DELETE TARGET;
      };
      CREATE REQUIRED PROPERTY molecular_number: std::str;
      CREATE LINK dissected_by: people::Person;
      CREATE REQUIRED PROPERTY morphological_code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
      ALTER PROPERTY morphological_code {
          CREATE REWRITE
              INSERT 
              USING ((((__subject__.biomat.code ++ '[') ++ .tube) ++ ']'));
          CREATE REWRITE
              UPDATE 
              USING ((((__subject__.biomat.code ++ '[') ++ .tube) ++ ']'));
      };
      CREATE ANNOTATION std::description := 'A single specimen isolated in a tube.';
      CREATE PROPERTY molecular_code: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE seq::ChromatoPrimer EXTENDING seq::Primer, default::Auditable;
  CREATE TYPE seq::SequencingInstitute EXTENDING default::Auditable {
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE SCALAR TYPE seq::ChromatoQuality EXTENDING enum<Contaminated, Failure, Ok, Unknown>;
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
  CREATE SCALAR TYPE seq::ExternalSeqType EXTENDING enum<NCBI, PERSCOM, LEGACY> {
      CREATE ANNOTATION std::description := "The sequence origin. 'PERSCOM' is 'Personal communication', 'Legacy' indicates the sequence originates from the lab but could not be registered as such due to missing required metadata.";
  };
  CREATE TYPE seq::ExternalSequence EXTENDING seq::Sequence, occurrence::Occurrence, default::Auditable {
      CREATE PROPERTY accession_number: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(12);
          CREATE CONSTRAINT std::min_len_value(6);
          CREATE ANNOTATION std::description := 'NCBI accession number of the original sequence.';
      };
      CREATE REQUIRED PROPERTY type: seq::ExternalSeqType;
      CREATE CONSTRAINT std::expression ON (EXISTS (.accession_number)) EXCEPT ((.type != seq::ExternalSeqType.NCBI));
      CREATE PROPERTY original_taxon: std::str {
          CREATE ANNOTATION std::description := 'The verbatim identification provided in the original source.';
      };
      CREATE REQUIRED PROPERTY specimen_identifier: std::str {
          CREATE ANNOTATION std::description := 'An identifier for the organism from which the sequence was produced, provided in the original source';
      };
  };
  CREATE SCALAR TYPE occurrence::QuantityType EXTENDING enum<Exact, Unknown, One, Several, Ten, Tens, Hundred>;
  CREATE TYPE occurrence::OccurrenceReport EXTENDING occurrence::Occurrence {
      CREATE LINK reference: reference::Article;
      CREATE LINK reported_by: people::Person;
      CREATE PROPERTY original_link: std::str;
      CREATE REQUIRED PROPERTY quantity: tuple<precision: occurrence::QuantityType, exact: std::int16> {
          CREATE CONSTRAINT std::expression ON ((((.precision = occurrence::QuantityType.Exact) AND (.exact > 0)) OR (.precision != occurrence::QuantityType.Exact)));
          CREATE REWRITE
              INSERT 
              USING (((
                  precision := __subject__.quantity.precision,
                  exact := -1
              ) IF (__subject__.quantity.precision != occurrence::QuantityType.Exact) ELSE __subject__.quantity));
          CREATE REWRITE
              UPDATE 
              USING (((
                  precision := __subject__.quantity.precision,
                  exact := -1
              ) IF (__subject__.quantity.precision != occurrence::QuantityType.Exact) ELSE __subject__.quantity));
      };
      CREATE PROPERTY specimen_available: tuple<collection: std::str, item: std::str> {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE seq::PCRSpecificity EXTENDING default::Auditable {
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE samples::BundledSpecimens EXTENDING samples::Sample {
      CREATE ANNOTATION std::description := 'A tube containing several specimens.';
      CREATE REQUIRED PROPERTY quantity: std::int16 {
          CREATE CONSTRAINT std::min_value(2);
      };
  };
  CREATE TYPE location::Locality EXTENDING default::Auditable {
      CREATE REQUIRED LINK country: location::Country;
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE ANNOTATION std::description := "Format like 'municipality|region[country_code]'";
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY municipality: std::str;
      CREATE PROPERTY region: std::str;
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING ((((((.municipality ?? 'NA') ++ '|') ++ (.region ?? 'NA')) ++ '|') ++ .country.code));
          CREATE REWRITE
              UPDATE 
              USING ((((((.municipality ?? 'NA') ++ '|') ++ (.region ?? 'NA')) ++ '|') ++ .country.code));
      };
      CREATE CONSTRAINT std::exclusive ON ((.region, .municipality));
  };
  CREATE TYPE event::EventDataset EXTENDING default::Auditable {
      CREATE MULTI LINK abiotic_measurements: event::AbioticMeasurement {
          ON TARGET DELETE ALLOW;
      };
      CREATE REQUIRED MULTI LINK maintainers: people::Person {
          ON TARGET DELETE RESTRICT;
      };
      CREATE MULTI LINK published_in: reference::Article {
          ON TARGET DELETE ALLOW;
      };
      CREATE MULTI LINK samplings: event::Sampling {
          ON TARGET DELETE ALLOW;
      };
      CREATE MULTI LINK spottings: event::Spotting {
          ON TARGET DELETE ALLOW;
      };
      CREATE PROPERTY comments: std::str;
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::max_len_value(40);
          CREATE CONSTRAINT std::min_len_value(4);
      };
  };
  CREATE TYPE seq::PCRReversePrimer EXTENDING seq::Primer, default::Auditable;
  ALTER TYPE event::Event {
      CREATE REQUIRED LINK site: location::Site;
      CREATE MULTI LINK programs: event::Program;
  };
  ALTER TYPE event::AbioticMeasurement {
      CREATE REQUIRED LINK param: event::AbioticParameter;
      CREATE CONSTRAINT std::exclusive ON ((.site, .param, .performed_on));
  };
  ALTER TYPE event::Sampling {
      CREATE MULTI LINK fixatives: default::Conservation;
      CREATE MULTI LINK measurements := (.<related_sampling[IS event::AbioticMeasurement]);
      CREATE MULTI LINK external_seqs := (.<sampling[IS seq::ExternalSequence]);
      CREATE MULTI LINK reports := (.<sampling[IS occurrence::OccurrenceReport]);
      CREATE MULTI LINK samples := (.<sampling[IS samples::BioMaterial]);
  };
  ALTER TYPE event::Spotting {
      CREATE MULTI LINK target_taxa: taxonomy::Taxon;
  };
  ALTER TYPE occurrence::Identification {
      CREATE REQUIRED LINK taxon: taxonomy::Taxon;
  };
  ALTER TYPE people::Person {
      CREATE MULTI LINK institutions: people::Institution;
      CREATE LINK user := (.<identity[IS people::User]);
      CREATE OPTIONAL PROPERTY role := (.user.role);
  };
  ALTER TYPE samples::BioMaterial {
      CREATE MULTI LINK specimens := (.<biomat[IS samples::Specimen]);
  };
  ALTER TYPE samples::Slide {
      CREATE REQUIRED LINK storage: storage::SlideStorage;
      CREATE CONSTRAINT std::exclusive ON ((.storage, .storage_position));
      CREATE REQUIRED LINK specimen: samples::Specimen;
  };
  ALTER TYPE seq::DNA {
      CREATE REQUIRED LINK specimen: samples::Specimen;
      CREATE MULTI LINK PCRs := (.<DNA[IS seq::PCR]);
      CREATE REQUIRED LINK extraction_method: seq::DNAExtractionMethod;
      CREATE REQUIRED LINK stored_in: storage::DNAStorage;
  };
  ALTER TYPE seq::PCR {
      CREATE MULTI LINK chromatograms := (.<PCR[IS seq::Chromatogram]);
      CREATE REQUIRED LINK forward_primer: seq::PCRForwardPrimer;
      CREATE REQUIRED LINK reverse_primer: seq::PCRReversePrimer;
  };
  CREATE TYPE admin::EmailSettings {
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
          SET default := 'DarCo';
          CREATE CONSTRAINT std::max_len_value(20);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE REQUIRED PROPERTY public: std::bool {
          SET default := true;
          CREATE ANNOTATION std::description := 'Whether parts of the platform are open to the public (anonymous users).';
      };
  };
  CREATE TYPE admin::SecuritySettings {
      CREATE REQUIRED PROPERTY account_token_lifetime: std::int32 {
          SET default := 24;
          CREATE ANNOTATION std::description := 'Validity period for an account operation token in hours';
          CREATE CONSTRAINT std::min_value(1);
      };
      CREATE REQUIRED PROPERTY auth_token_lifetime: std::int32 {
          SET default := 60;
          CREATE ANNOTATION std::description := 'Validity period for an authentication token in seconds';
          CREATE CONSTRAINT std::min_value(60);
      };
      CREATE REQUIRED PROPERTY jwt_secret_key: std::str;
      CREATE REQUIRED PROPERTY min_password_strength: std::int32 {
          SET default := 3;
          CREATE CONSTRAINT std::max_value(5);
          CREATE CONSTRAINT std::min_value(3);
      };
  };
  CREATE TYPE admin::Settings {
      CREATE LINK email: admin::EmailSettings {
          ON SOURCE DELETE DELETE TARGET;
      };
      CREATE REQUIRED LINK instance: admin::InstanceSettings {
          ON SOURCE DELETE DELETE TARGET;
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED LINK security: admin::SecuritySettings {
          ON SOURCE DELETE DELETE TARGET;
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE datasets::MOTU {
      CREATE REQUIRED LINK dataset: datasets::MOTUDataset;
  };
  ALTER TYPE datasets::MOTUDataset {
      CREATE MULTI LINK MOTUs := (.<dataset[IS datasets::MOTU]);
  };
  ALTER TYPE location::Site {
      CREATE MULTI LINK abiotic_measurements := (.<site[IS event::AbioticMeasurement]);
      CREATE MULTI LINK samplings := (.<site[IS event::Sampling]);
      CREATE MULTI LINK spottings := (.<site[IS event::Spotting]);
      CREATE REQUIRED LINK locality: location::Locality {
          ON TARGET DELETE RESTRICT;
      };
  };
  ALTER TYPE seq::ExternalSequence {
      CREATE LINK source_sample: occurrence::OccurrenceReport;
  };
  ALTER TYPE occurrence::OccurrenceReport {
      CREATE MULTI LINK sequences := (.<source_sample[IS seq::ExternalSequence]);
  };
  ALTER TYPE seq::AssembledSequence {
      CREATE REQUIRED LINK specimen: samples::Specimen;
      CREATE CONSTRAINT std::exclusive ON ((.specimen, .is_reference)) EXCEPT (NOT (.is_reference));
      CREATE REQUIRED MULTI LINK chromatograms: seq::Chromatogram;
  };
  ALTER TYPE samples::Specimen {
      CREATE MULTI LINK sequences := (.<specimen[IS seq::AssembledSequence]);
      CREATE LINK identification := (((SELECT
          .sequences
      FILTER
          .is_reference
      )).identification);
      ALTER PROPERTY molecular_number {
          CREATE REWRITE
              INSERT 
              USING (<std::str>(SELECT
                  std::count(samples::Specimen)
              FILTER
                  (samples::Specimen.biomat.sampling.site = __subject__.biomat.sampling.site)
              ));
      };
      CREATE MULTI LINK slides := (.<specimen[IS samples::Slide]);
  };
  ALTER TYPE samples::BioMaterial {
      CREATE MULTI LINK identified_taxa := (SELECT
          (DISTINCT (.specimens.identification.taxon) ?? .identification.taxon)
      );
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (((.identification.taxon.code ++ '|') ++ .sampling.code));
          CREATE REWRITE
              UPDATE 
              USING (((.identification.taxon.code ++ '|') ++ .sampling.code));
      };
      CREATE MULTI LINK published_in: reference::Article;
      CREATE MULTI LINK bundles := (.<biomat[IS samples::BundledSpecimens]);
      CREATE MULTI LINK content := (.<biomat[IS samples::Sample]);
  };
  ALTER TYPE event::Sampling {
      CREATE MULTI LINK occurring_taxa := (WITH
          ext_samples_no_seqs := 
              (SELECT
                  .reports
              FILTER
                  NOT (EXISTS (.sequences))
              )
      SELECT
          DISTINCT (((ext_samples_no_seqs.identification.taxon UNION .external_seqs.identification.taxon) UNION .samples.identified_taxa))
      );
      CREATE MULTI LINK methods: event::SamplingMethod;
      CREATE PROPERTY generated_code := (SELECT
          (((.site.code ++ '_') ++ <std::str>std::datetime_get(.performed_on.date, 'year')) ++ <std::str>std::datetime_get(.performed_on.date, 'month'))
      );
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (SELECT
                  (((.generated_code ++ '_') ++ <std::str>(SELECT
                      std::count(event::Sampling)
                  FILTER
                      (event::Sampling.code = __subject__.generated_code)
                  )) IF (SELECT
                      EXISTS (event::Sampling)
                  FILTER
                      (event::Sampling.code = __subject__.generated_code)
                  ) ELSE .generated_code)
              );
      };
      CREATE MULTI LINK target_taxa: taxonomy::Taxon {
          CREATE ANNOTATION std::title := 'Taxonomic groups that were the target of the sampling effort';
      };
  };
  ALTER TYPE location::Country {
      CREATE MULTI LINK localities := (.<country[IS location::Locality]);
  };
  ALTER TYPE location::Locality {
      CREATE MULTI LINK sites := (.<locality[IS location::Site]);
  };
  ALTER TYPE location::SiteDataset {
      CREATE MULTI LINK sites: location::Site;
  };
  CREATE ABSTRACT TYPE people::AccountToken {
      CREATE REQUIRED LINK user: people::User {
          ON TARGET DELETE DELETE SOURCE;
          CREATE DELEGATED CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY expires: std::datetime;
      CREATE REQUIRED PROPERTY token: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE people::EmailConfirmation EXTENDING people::AccountToken;
  CREATE TYPE people::PasswordReset EXTENDING people::AccountToken;
  ALTER TYPE people::Institution {
      CREATE MULTI LINK people := (.<institutions[IS people::Person]);
  };
  CREATE TYPE people::PendingUserRequest {
      CREATE REQUIRED LINK user: people::User {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE DELETE SOURCE;
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY created_on: std::datetime {
          CREATE REWRITE
              INSERT 
              USING (std::datetime_of_statement());
      };
      CREATE REQUIRED PROPERTY identity: tuple<first_name: std::str, last_name: std::str, institution: std::str>;
      CREATE REQUIRED PROPERTY motive: std::str;
  };
  CREATE TYPE people::UserInvitation {
      CREATE REQUIRED LINK identity: people::Person;
      CREATE REQUIRED PROPERTY expires: std::datetime;
      CREATE REQUIRED PROPERTY role: people::UserRole;
      CREATE REQUIRED PROPERTY token: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
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
