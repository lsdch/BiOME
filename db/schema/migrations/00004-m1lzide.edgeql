CREATE MIGRATION m1lzidennerhb2nvvhipwh4ecvwhdgimsdx45v7vcpiqizl6nyshja
    ONTO m1tswrc4zoyyku5yqoavbptqodgu36bc4bfjw2s4gzl34g7jd6vr2a
{
  ALTER TYPE default::Vocabulary {
      ALTER PROPERTY code {
          CREATE CONSTRAINT std::max_len_value(12);
      };
  };
  ALTER TYPE default::Vocabulary {
      ALTER PROPERTY code {
          DROP CONSTRAINT std::max_len_value(8);
      };
  };
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
              USING ((IF (SELECT
                  EXISTS (events::Sampling)
              FILTER
                  (events::Sampling.code = __subject__.generated_code)
              ) THEN (SELECT
                  ((.generated_code ++ '_') ++ <std::str>(SELECT
                      std::count(events::Sampling)
                  FILTER
                      (events::Sampling.code = __subject__.generated_code)
                  ))
              ) ELSE .generated_code));
      };
      ALTER PROPERTY is_donation {
          SET default := false;
      };
  };
};
