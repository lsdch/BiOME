CREATE MIGRATION m1bu4c6lntjbgvenwasslcmnjzsordq233weqlo7d4mkjmze2hfurq
    ONTO m1qwlpsc5dcohpkmplks5s7xbhrlxoklqlyonpzxd44wnp3sagkd5a
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
              USING ((IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE std::str_replace(.name, ' ', '_')));
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
              USING ((IF (__specified__.code AND (std::len(.code) > 0)) THEN .code ELSE std::str_replace(.name, ' ', '_')));
      };
  };
  ALTER SCALAR TYPE taxonomy::TaxonStatus EXTENDING enum<Accepted, Unreferenced, Unclassified>;
};
