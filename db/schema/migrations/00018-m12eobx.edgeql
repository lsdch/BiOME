CREATE MIGRATION m12eobx6drp6r2akldqqg24s5wadmhrnxgrtiztdj3lu2uhrkbknvq
    ONTO m15srydi2fdx3ra3hwsnjy7qtukvdzt75i665jjsfecav2efj3ehiq
{
  ALTER FUNCTION references::generate_article_code(authors: array<std::str>, year: std::int32) USING (SELECT
      (IF (std::len(authors) = 1) THEN (((std::str_split((authors)[0], ' '))[0] ++ '_') ++ <std::str>year) ELSE (IF (std::len(authors) = 2) THEN ((std::array_join([(std::str_split((authors)[0], ' '))[0], (std::str_split((authors)[1], ' '))[0]], '_') ++ '_') ++ <std::str>year) ELSE (((std::str_split((authors)[0], ' '))[0] ++ '_et_al_') ++ <std::str>year)))
  );
};
