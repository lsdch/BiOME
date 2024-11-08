CREATE MIGRATION m1eoedhqazvntxpkwxdxweink5jmtpdvnnahsc4sriftlfxkb35kyq
    ONTO m1w7dosbzbahixus6ar3lkxxtnxxu7uunqnpaveyn2uqg4w7qrioma
{
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (SELECT
                  ((.generated_code ++ '.') ++ <std::str>(SELECT
                      std::count(DETACHED events::Sampling FILTER
                          (.generated_code = __subject__.generated_code)
                      )
                  ))
              );
      };
  };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  ((.generated_code ++ '.') ++ <std::str>(SELECT
                      std::count(DETACHED events::Sampling FILTER
                          (.generated_code = __subject__.generated_code)
                      )
                  ))
              );
      };
  };
};
