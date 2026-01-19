CREATE MIGRATION m1yljweecpms65j2ngtg6cgk2tihbf5bm72nukz7uas5kgvjjr5tjq
    ONTO m1sfuoo3wfnz2dbeb7ycflmo3pz4rwbwvgwelur4k75e2nvfdjeowa
{
  DROP ALIAS seq::GenericSequenceWithDetails;
  ALTER TYPE datasets::MOTUDataset {
      DROP LINK generated_by;
  };
  ALTER TYPE datasets::MOTUDataset {
      CREATE REQUIRED MULTI PROPERTY generated_by: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE datasets::ResearchProgram {
      DROP LINK managers;
  };
  ALTER TYPE datasets::ResearchProgram {
      CREATE REQUIRED MULTI PROPERTY managers: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE events::Action {
      DROP LINK performed_by;
  };
  ALTER TYPE events::Action {
      CREATE MULTI PROPERTY performed_by: std::str;
      DROP LINK performed_by_groups;
  };
  ALTER TYPE occurrence::Identification {
      DROP LINK identified_by;
  };
  ALTER TYPE occurrence::Identification {
      CREATE MULTI PROPERTY identified_by: std::str;
  };
  ALTER TYPE references::DataSource {
      CREATE PROPERTY contact: std::str;
  };
  ALTER TYPE samples::Slide {
      DROP LINK mounted_by;
  };
  ALTER TYPE samples::Slide {
      CREATE REQUIRED MULTI PROPERTY mounted_by: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE samples::Specimen {
      DROP LINK dissected_by;
  };
  ALTER TYPE samples::Specimen {
      CREATE MULTI PROPERTY dissected_by: std::str;
  };
  ALTER TYPE seq::AssembledSequence {
      DROP LINK assembled_by;
  };
  ALTER TYPE seq::AssembledSequence {
      CREATE REQUIRED MULTI PROPERTY assembled_by: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE sequencing::BatchRequest {
      DROP LINK requested_by;
  };
  ALTER TYPE sequencing::BatchRequest {
      DROP LINK requested_to;
  };
  ALTER TYPE sequencing::BatchRequest {
      CREATE REQUIRED MULTI PROPERTY requested_by: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE sequencing::BatchRequest {
      CREATE MULTI PROPERTY requested_to: std::str;
  };
  ALTER TYPE sequencing::DNA {
      DROP LINK extracted_by;
  };
  ALTER TYPE sequencing::DNA {
      CREATE REQUIRED MULTI PROPERTY extracted_by: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE storage::Collection {
      DROP LINK maintainers;
  };
  ALTER TYPE storage::Collection {
      CREATE REQUIRED MULTI PROPERTY maintainers: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
};
