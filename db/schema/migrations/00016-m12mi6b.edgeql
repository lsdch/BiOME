CREATE MIGRATION m12mi6bd2jaqj2e66grrosjljycuazdysn44h3s7aw7x32hbwvvphq
    ONTO m16pc4ljrpwskombxhanj7mmzsyg4zbcpipdtdutykscb35utqnpwq
{
  ALTER TYPE references::Article {
      ALTER PROPERTY code {
          SET default := (WITH
              gen_code := 
                  references::generate_article_code(.authors, .year)
              ,
              discriminant := 
                  (SELECT
                      (std::count(DETACHED references::Article FILTER
                          (references::generate_article_code(.authors, .year) = gen_code)
                      ) - 1)
                  )
          SELECT
              (IF (discriminant >= 0) THEN (references::generate_article_code(.authors, .year) ++ (GLOBAL references::alphabet)[discriminant]) ELSE references::generate_article_code(.authors, .year))
          );
      };
  };
};
