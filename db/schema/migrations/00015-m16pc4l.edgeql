CREATE MIGRATION m16pc4ljrpwskombxhanj7mmzsyg4zbcpipdtdutykscb35utqnpwq
    ONTO m1qynlncy35fdbvsj5otxf6fk5txcpazm2wlgk6ov7paujmhantfyq
{
  ALTER TYPE references::Article {
      ALTER PROPERTY code {
          SET default := (WITH
              gen_code := 
                  references::generate_article_code(.authors, .year)
              ,
              discriminant := 
                  (SELECT
                      std::count(DETACHED references::Article FILTER
                          (references::generate_article_code(references::Article.authors, references::Article.year) = gen_code)
                      )
                  )
          SELECT
              (IF (discriminant > 0) THEN (references::generate_article_code(.authors, .year) ++ (GLOBAL references::alphabet)[(discriminant - 1)]) ELSE references::generate_article_code(.authors, .year))
          );
          DROP REWRITE
              INSERT ;
              DROP REWRITE
                  UPDATE ;
              };
          };
};
