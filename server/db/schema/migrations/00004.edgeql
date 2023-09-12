CREATE MIGRATION m1yfbgj6u2atxxlxngmizkeoi4wy3wkzmbz5yebiomi7abstldr5fa
    ONTO m1gwdc3fyv5f5evagrf6nmvjjuxy54mwqygepqyzegij5ob2idxrbq
{
  ALTER TYPE event::Sampling {
      CREATE MULTI LINK external_seqs := (.<sampling[IS seq::ExternalSequence]);
  };
  ALTER TYPE seq::ExternalSequence {
      CREATE LINK source_sample: occurrence::OccurrenceReport;
  };
  ALTER TYPE occurrence::OccurrenceReport {
      CREATE MULTI LINK sequences := (.<source_sample[IS seq::ExternalSequence]);
  };
  ALTER TYPE occurrence::Occurrence {
      DROP LINK identifications;
  };
  ALTER TYPE samples::BioMaterial {
      CREATE MULTI LINK specimens := (.<biomat[IS samples::Specimen]);
  };
  ALTER TYPE samples::Specimen {
      CREATE MULTI LINK sequences := (.<specimen[IS seq::AssembledSequence]);
  };
  ALTER TYPE samples::Specimen {
      CREATE LINK identification := (((SELECT
          .sequences
      FILTER
          .is_reference
      )).identification);
  };
  ALTER TYPE samples::BioMaterial {
      CREATE MULTI LINK identifications := (.specimens.identification);
      CREATE MULTI LINK bundles := (.<biomat[IS samples::BundledSpecimens]);
  };
  ALTER TYPE event::Sampling {
      CREATE MULTI LINK all_ids := (WITH
          ext_samples_no_seqs := 
              (SELECT
                  .reports
              FILTER
                  NOT (EXISTS (.sequences))
              )
      SELECT
          ((ext_samples_no_seqs.identification UNION .external_seqs.identification) UNION .samples.identifications)
      );
  };
  ALTER TYPE samples::Specimen {
      DROP LINK molecular_identification;
  };
};
