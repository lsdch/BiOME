CREATE MIGRATION m1dyelrc2hzqdh2bzhodrmvtqgt5bfaaohngico7wajkkvspb4zh3a
    ONTO m1fjy4i5edtczsyuksemhlucpuz6nxrj2v5d7r6rd7xwyhaxlhzfpq
{
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
