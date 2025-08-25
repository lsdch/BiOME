CREATE MIGRATION m1u6cxxvdqc7f3qxmjmt6j7dpirx3yg5dbgiwncm4x35zqu4w4lt5q
    ONTO m146tlry5a5gr32xrnqweiyozr4mhzrws2ojgfwly3kacqb4iijvaq
{
  CREATE FUNCTION taxonomy::taxon_code(taxon: taxonomy::Taxon) ->  std::str USING (std::re_replace(taxon.name, '[^()a-zA-Z]', '_'));
  ALTER FUNCTION occurrence::biomat_code(taxon: taxonomy::Taxon, sampling: events::Sampling) USING ((((taxonomy::taxon_code(taxon) ++ '[') ++ sampling.code) ++ ']'));
  ALTER FUNCTION seq::generate_ext_seq_code(seq: seq::ExternalSequence) USING (WITH
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
      ((((((taxonomy::taxon_code(seq.identification.taxon) ++ '[') ++ seq.sampling.code) ++ ']') ++ seq.specimen_identifier) ++ '|') ++ suffix)
  );
  DROP FUNCTION taxonomy::find_taxon(name_or_code: std::str);
  CREATE FUNCTION taxonomy::find_taxon(name: std::str) ->  taxonomy::Taxon USING (SELECT
      std::assert_single(std::assert_exists(taxonomy::Taxon FILTER
          (.name = name)
      , message := ('Failed to find taxon with name: ' ++ name)), message := ('Multiple taxa matching name: ' ++ name))
  );
  ALTER TYPE taxonomy::Taxon {
      CREATE INDEX ON ((.name, .rank, .status));
  };
  ALTER TYPE taxonomy::Taxon {
      DROP INDEX ON ((.name, .code, .rank, .status));
      DROP PROPERTY code;
  };
};
