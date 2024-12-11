CREATE MIGRATION m15srydi2fdx3ra3hwsnjy7qtukvdzt75i665jjsfecav2efj3ehiq
    ONTO m12mi6bd2jaqj2e66grrosjljycuazdysn44h3s7aw7x32hbwvvphq
{
  ALTER FUNCTION references::generate_article_code(authors: array<std::str>, year: std::int32) USING (SELECT
      (IF (std::count(authors) = 1) THEN (((std::str_split((authors)[0], ' '))[0] ++ '_') ++ <std::str>year) ELSE (IF (std::count(authors) = 2) THEN ((std::array_join([(std::str_split((authors)[0], ' '))[0], (std::str_split((authors)[1], ' '))[0]], '_') ++ '_') ++ <std::str>year) ELSE (((std::str_split((authors)[0], ' '))[0] ++ '_et_al_') ++ <std::str>year)))
  );
  ALTER TYPE occurrence::BioMaterial {
      ALTER LINK published_in {
          CREATE PROPERTY original_source: std::bool;
      };
  };
};
