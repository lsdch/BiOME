CREATE MIGRATION m1fvea757cwbtb22qn55tjz5n5z66looe3d5avd2licrlhsq4af74q
    ONTO m1wcpeg6y4xz6lbdolthika34bfm4xldvk4gkxl7u7jtgz7cr4naia
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
                  (IF (SELECT
                      EXISTS (events::Sampling)
                  FILTER
                      (events::Sampling.code = __subject__.generated_code)
                  ) THEN (SELECT
                      ((.generated_code ++ '.') ++ <std::str>(SELECT
                          std::count(events::Sampling)
                      FILTER
                          (events::Sampling.code = __subject__.generated_code)
                      ))
                  ) ELSE (SELECT
                      .generated_code
                  ))
              );
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  (IF (SELECT
                      EXISTS (events::Sampling)
                  FILTER
                      (events::Sampling.code = __subject__.generated_code)
                  ) THEN (SELECT
                      ((.generated_code ++ '.') ++ <std::str>(SELECT
                          std::count(events::Sampling)
                      FILTER
                          (events::Sampling.code = __subject__.generated_code)
                      ))
                  ) ELSE (SELECT
                      .generated_code
                  ))
              );
      };
  };
};
