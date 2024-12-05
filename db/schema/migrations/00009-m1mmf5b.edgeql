CREATE MIGRATION m1mmf5baxwd33am7pjwddrznuy7k6gral3yyhf5iljanbpuiupkyca
    ONTO m14wjjj7d46idowjhavxow7ipxz5pfnkklo4xms2fn2u4byipcm6jq
{
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
};
