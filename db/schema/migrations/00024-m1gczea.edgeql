CREATE MIGRATION m1gczeafpaw6yy7f6x6jwlkvu7q42r4vvzd3va2w54k54bolm37tcq
    ONTO m1okr4aq5ymvtsyoxg7luhjhvhoepaw45hb5jkemoc6kh7m6kpfqwa
{
  CREATE TYPE seq::ExternalSeqOrigin EXTENDING default::Auditable, default::Vocabulary {
      CREATE REQUIRED PROPERTY accession_required: std::bool {
          SET default := false;
      };
      CREATE PROPERTY link_template: std::str;
  };
  ALTER TYPE seq::ExternalSequence {
      CREATE REQUIRED LINK origin: seq::ExternalSeqOrigin {
          SET REQUIRED USING (<seq::ExternalSeqOrigin>{});
      };
  };
  ALTER TYPE seq::ExternalSequence {
      DROP PROPERTY type;
  };
  DROP SCALAR TYPE seq::ExternalSeqType;
};
