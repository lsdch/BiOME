CREATE MIGRATION m1fgjgxn27luncdge6xbyn4d3i3mpf3ita2ue4kpls2nr6fupf53ja
    ONTO m1mmf5baxwd33am7pjwddrznuy7k6gral3yyhf5iljanbpuiupkyca
{
  ALTER TYPE events::Sampling {
      DROP LINK occurring_taxa;
  };
  ALTER TYPE occurrence::InternalBioMat {
      DROP LINK identified_taxa;
  };
  ALTER TYPE samples::Specimen {
      DROP LINK identification;
  };
  ALTER TYPE seq::AssembledSequence {
      DROP LINK identification;
  };
  ALTER TYPE seq::AssembledSequence {
      CREATE LINK identification: occurrence::Identification {
          ON SOURCE DELETE DELETE TARGET;
          SET REQUIRED USING (<occurrence::Identification>{});
      };
      CREATE REQUIRED LINK sampling: events::Sampling {
          SET REQUIRED USING (<events::Sampling>{});
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
      DROP EXTENDING default::Auditable;
      EXTENDING occurrence::Occurrence LAST;
  };
  ALTER TYPE seq::AssembledSequence {
      ALTER LINK identification {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
  ALTER TYPE samples::Specimen {
      CREATE LINK identification := (((SELECT
          .sequences
      FILTER
          .is_reference
      )).identification);
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE MULTI LINK identified_taxa := (SELECT
          (DISTINCT (.specimens.identification.taxon) ?? .identification.taxon)
      );
  };
  ALTER TYPE seq::ExternalSequence {
      DROP CONSTRAINT std::expression ON (EXISTS (.accession_number)) EXCEPT ((.type != seq::ExternalSeqType.NCBI));
      DROP EXTENDING default::Auditable;
  };
  ALTER TYPE events::Sampling {
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
  ALTER TYPE occurrence::BioMaterial {
      ALTER PROPERTY code {
          ALTER ANNOTATION std::description := "Format like 'taxon_short_code|sampling_code'";
      };
  };
};
