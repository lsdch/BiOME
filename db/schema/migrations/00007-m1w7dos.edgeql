CREATE MIGRATION m1w7dosbzbahixus6ar3lkxxtnxxu7uunqnpaveyn2uqg4w7qrioma
    ONTO m1fvea757cwbtb22qn55tjz5n5z66looe3d5avd2licrlhsq4af74q
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
                  (IF EXISTS ((SELECT
                      DETACHED events::Sampling
                  FILTER
                      (.generated_code = __subject__.generated_code)
                  )) THEN (SELECT
                      ((.generated_code ++ '.') ++ <std::str>(SELECT
                          std::count(DETACHED events::Sampling FILTER
                              (.generated_code = __subject__.generated_code)
                          )
                      ))
                  ) ELSE (SELECT
                      .generated_code
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
                  (IF EXISTS ((SELECT
                      DETACHED events::Sampling
                  FILTER
                      (.generated_code = __subject__.generated_code)
                  )) THEN (SELECT
                      ((.generated_code ++ '.') ++ <std::str>(SELECT
                          std::count(DETACHED events::Sampling FILTER
                              (.generated_code = __subject__.generated_code)
                          )
                      ))
                  ) ELSE (SELECT
                      .generated_code
                  ))
              );
      };
  };
};
