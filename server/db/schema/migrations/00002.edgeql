CREATE MIGRATION m1nq6igapvk5az34hdeuik4vgqwy6kfgra73qleshkrzrmedjfkzlq
    ONTO m1lmvbdy3bphrpfe4boazmyzvabpv23iphgudaxxwbnoo5gzc6u32a
{
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY code {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (WITH
                  chopped := 
                      std::str_split(.name, ' ')
                  ,
                  suffix := 
                      ('[syn]' IF (.status = taxonomy::TaxonStatus.Synonym) ELSE '')
              SELECT
                  ((.name ++ suffix) IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix))
              );
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (WITH
                  chopped := 
                      std::str_split(.name, ' ')
                  ,
                  suffix := 
                      ('[syn]' IF (.status = taxonomy::TaxonStatus.Synonym) ELSE '')
              SELECT
                  ((.name ++ suffix) IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix))
              );
      };
  };
};
