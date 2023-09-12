CREATE MIGRATION m1zekcs5zvbhydwi6zxeliwrmbwozwwuef3e2bsqsvtwwmvroxkfuq
    ONTO m1yfbgj6u2atxxlxngmizkeoi4wy3wkzmbz5yebiomi7abstldr5fa
{
  ALTER TYPE event::Sampling {
      ALTER LINK all_ids {
          USING (WITH
              ext_samples_no_seqs := 
                  (SELECT
                      .reports
                  FILTER
                      NOT (EXISTS (.sequences))
                  )
          SELECT
              ((DISTINCT (ext_samples_no_seqs.identification.taxon) UNION .external_seqs.identification.taxon) UNION .samples.identifications.taxon)
          );
      };
  };
};
