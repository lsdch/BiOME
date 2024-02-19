CREATE MIGRATION m1bzpahvfvejjec2edpihodpfdzkj4bswrnlhr2x42sgsw4f3apypa
    ONTO m1r7yjx2icbvr7rj6c5oka2dmntkeakg7bfgws2yinf6d5ftgl5lva
{
  ALTER TYPE people::User {
      ALTER PROPERTY password {
          CREATE REWRITE
              INSERT 
              USING (ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des')));
          CREATE REWRITE
              UPDATE 
              USING (ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des')));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY GBIF_ID {
          RESET OPTIONALITY;
      };
  };
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
                  (.code IF (EXISTS (__specified__.code) AND __specified__.code) ELSE ((.name ++ suffix) IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix)))
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
                  (.code IF (EXISTS (__specified__.code) AND __specified__.code) ELSE ((.name ++ suffix) IF NOT ((.rank IN {taxonomy::Rank.Species, taxonomy::Rank.Subspecies})) ELSE ((std::str_upper((chopped)[0][:3]) ++ std::array_join((chopped)[1:], '_')) ++ suffix)))
              );
      };
  };
};
