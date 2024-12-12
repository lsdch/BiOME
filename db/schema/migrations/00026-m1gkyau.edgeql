CREATE MIGRATION m1gkyaucljimawntjiae6426sjzubgnes5m5jfvoiwiqu6lld27grq
    ONTO m1x4b3ztm5kjvp5nsksrprwplv2sjbg43bn47huk54tlwe7zqffroq
{
  ALTER TYPE seq::ExternalSeqOrigin {
      DROP PROPERTY accession_required;
      DROP PROPERTY link_template;
  };
  ALTER TYPE seq::ExternalSequence {
      DROP LINK origin;
  };
  DROP TYPE seq::ExternalSeqOrigin;
  ALTER TYPE seq::ExternalSequence {
      ALTER LINK reference {
          RENAME TO published_in;
      };
  };
  CREATE TYPE seq::SeqDB EXTENDING default::Auditable, default::Vocabulary {
      CREATE PROPERTY link_template: std::str;
  };
  CREATE TYPE seq::SeqReference {
      CREATE REQUIRED LINK db: seq::SeqDB {
          ON TARGET DELETE DELETE SOURCE;
      };
      CREATE REQUIRED PROPERTY accession: std::str {
          CREATE CONSTRAINT std::max_len_value(32);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE CONSTRAINT std::exclusive ON ((.db, .accession));
  };
  ALTER TYPE seq::ExternalSequence {
      CREATE MULTI LINK referenced_in: seq::SeqReference {
          CREATE CONSTRAINT std::exclusive;
      };
      DROP PROPERTY accession_number;
  };
  CREATE SCALAR TYPE seq::ExtSeqOrigin EXTENDING enum<Lab, PersCom, DB>;
  ALTER TYPE seq::ExternalSequence {
      CREATE REQUIRED PROPERTY origin: seq::ExtSeqOrigin {
          SET REQUIRED USING (<seq::ExtSeqOrigin>{});
      };
  };
};
