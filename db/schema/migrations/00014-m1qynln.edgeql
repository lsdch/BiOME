CREATE MIGRATION m1qynlncy35fdbvsj5otxf6fk5txcpazm2wlgk6ov7paujmhantfyq
    ONTO m1muksybuie6vuxtspjeydwd6esqd2tzgws2yeh5vlrvdzn76tkiva
{
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          sampling: {
              *
          },
          identification: {
              *
          },
          meta: {
              *
          },
          category := (IF (occurrence::BioMaterial IS occurrence::InternalBioMat) THEN 'Internal' ELSE (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN 'External' ELSE 'Unknown'))
      }
  );
  CREATE FUNCTION references::generate_article_code(authors: array<std::str>, year: std::int32) ->  std::str USING (SELECT
      (IF (std::count(authors) = 1) THEN ((std::str_split((authors)[0], ' '))[0] ++ <std::str>year) ELSE (IF (std::count(authors) = 2) THEN (std::array_join([(std::str_split((authors)[0], ' '))[0], (std::str_split((authors)[1], ' '))[0]], '_') ++ <std::str>year) ELSE (((std::str_split((authors)[0], ' '))[0] ++ '_et_al_') ++ <std::str>year)))
  );
  ALTER TYPE references::Article {
      ALTER PROPERTY authors {
          RESET CARDINALITY USING (SELECT
              .authors 
          LIMIT
              1
          );
          SET REQUIRED USING (<std::str>{});
          SET TYPE array<std::str> USING (std::array_agg(.authors));
      };
      CREATE REQUIRED PROPERTY code: std::str {
          SET default := (references::generate_article_code(.authors, .year));
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE GLOBAL references::alphabet := ('abcdefghijklmnopqrstuvxyz');
  ALTER TYPE references::Article {
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (WITH
                  discriminant := 
                      (SELECT
                          std::count(DETACHED references::Article FILTER
                              (references::Article.code = __subject__.code)
                          )
                      )
              SELECT
                  (IF (discriminant > 0) THEN (references::generate_article_code(.authors, .year) ++ (GLOBAL references::alphabet)[(discriminant - 1)]) ELSE references::generate_article_code(.authors, .year))
              );
          CREATE REWRITE
              UPDATE 
              USING (WITH
                  discriminant := 
                      (SELECT
                          std::count(DETACHED references::Article FILTER
                              (references::Article.code = __subject__.code)
                          )
                      )
              SELECT
                  (IF (discriminant > 0) THEN (references::generate_article_code(.authors, .year) ++ (GLOBAL references::alphabet)[(discriminant - 1)]) ELSE references::generate_article_code(.authors, .year))
              );
      };
      CREATE INDEX ON ((.code, .year));
      ALTER PROPERTY title {
          RESET OPTIONALITY;
      };
  };
  ALTER TYPE occurrence::InternalBioMat {
      DROP LINK content;
  };
};
