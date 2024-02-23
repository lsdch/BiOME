CREATE MIGRATION m1ejkaogorsq45wgejygvtptjx4zapaks7sjuwk5c775nml6fcalga
    ONTO m1dv3nnnmq6l2zj64v7pn2ptqvk4rz2uckzlusijvy7zj7dgkrk6fa
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
                  (IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE (IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) THEN (.name ++ suffix) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix)))
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
                  (IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE (IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) THEN (.name ++ suffix) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix)))
              );
      };
  };
};
