CREATE MIGRATION m14o4u2w6uykl5oido6eamuoamtgtmhs74k56nj77qf57ge7tgtiiq
    ONTO m1fuslhbl7c67z5afycwwiupsuix2c5k77aoxf2svyju64d25a5bsq
{
  CREATE FUNCTION seq::generate_ext_seq_code(seq: seq::ExternalSequence) ->  std::str USING (WITH
      suffix := 
          (IF (seq.origin = seq::ExtSeqOrigin.Lab) THEN 'lab' ELSE (IF (seq.origin = seq::ExtSeqOrigin.PersCom) THEN 'perscom' ELSE (WITH
              sources := 
                  (SELECT
                      seq.referenced_in
                  FILTER
                      seq.referenced_in.is_origin
                  )
          SELECT
              std::array_join(std::array_agg(sources.code), '|')
          )))
  SELECT
      ((((((seq.identification.taxon.code ++ '[') ++ seq.sampling.code) ++ ']') ++ seq.specimen_identifier) ++ '|') ++ suffix)
  );
  ALTER TYPE seq::ExternalSequence {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE seq::ExternalSequence {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  seq::generate_ext_seq_code(__subject__)
              );
      };
  };
  ALTER TYPE settings::AbstractSettingsSpec {
      CREATE PROPERTY description: std::str;
  };
};
