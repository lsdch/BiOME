CREATE MIGRATION m1nkgduvf2df6eqyrsimwsqmfpsbeu25obmxcsy5klvv2eedzjxvsa
    ONTO m1ha3vu5s6nlivik2ecyyaomjpewtqlfvo5e575uxzp7fqurupszma
{
  DROP FUNCTION occurrence::insert_external_occurrence(sampling: events::Sampling, data: std::json);
  DROP FUNCTION occurrence::insert_internal_biomat(sampling: events::Sampling, data: std::json);
  ALTER TYPE occurrence::Occurrence {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  DROP FUNCTION occurrence::occurrence_code(taxon: taxonomy::Taxon, sampling: events::Sampling);
  CREATE FUNCTION occurrence::occurrence_code(taxon: taxonomy::Taxon, sampling_code: std::str) ->  std::str USING ((((taxonomy::taxon_code(taxon) ++ '[') ++ sampling_code) ++ ']'));
  ALTER TYPE events::Sampling {
      CREATE REQUIRED PROPERTY code := (((((.site.code ++ '.') ++ <std::str>.number) ++ '|') ++ date::to_code(.performed_on)));
  };
  ALTER TYPE occurrence::Occurrence {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.code THEN .code ELSE occurrence::occurrence_code(.identification.taxon, .sampling.code)));
      };
  };
};
