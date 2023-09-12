CREATE MIGRATION m1y5ynh7hph447hhf2dmz6lpi5lkaj2pnueuxk3qqr3pvglilephxq
    ONTO m1c2b4m5ojl5knshpexhq5dnc3gi7hklyszcium4ddbrvi6cg35vkq
{
  ALTER TYPE samples::BioMaterial {
      CREATE MULTI LINK identified_taxa := (SELECT
          (DISTINCT (.specimens.identification.taxon) ?? .identification.taxon)
      );
  };
  ALTER TYPE event::Sampling {
      ALTER LINK occurring_taxa {
          USING (WITH
              ext_samples_no_seqs := 
                  (SELECT
                      .reports
                  FILTER
                      NOT (EXISTS (.sequences))
                  )
          SELECT
              ((DISTINCT (ext_samples_no_seqs.identification.taxon) UNION .external_seqs.identification.taxon) UNION .samples.identified_taxa)
          );
      };
  };
  ALTER TYPE samples::BioMaterial {
      DROP LINK identifications;
  };
};
