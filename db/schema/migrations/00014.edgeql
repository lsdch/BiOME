CREATE MIGRATION m1dv3nnnmq6l2zj64v7pn2ptqvk4rz2uckzlusijvy7zj7dgkrk6fa
    ONTO m1ur3upk4dbze2w4jahujyhe6pfhsdqdsz2yyf3gnwhajh353omnja
{
  CREATE FUNCTION people::middle_names_initials(s: OPTIONAL std::str) ->  std::str USING ((IF EXISTS (s) THEN (WITH
      split_middle_names := 
          std::str_split(s, ' ')
  SELECT
      std::array_join(std::array_agg((FOR middle_name IN std::array_unpack(split_middle_names)
      UNION 
          (middle_name)[0])), '')
  ) ELSE (SELECT
      ''
  )));
  ALTER TYPE people::Person {
      CREATE REQUIRED PROPERTY alias: std::str {
          SET default := (<std::str>{});
          CREATE REWRITE
              INSERT 
              USING ((default::null_if_empty(.alias) ?? (WITH
                  middle_names := 
                      people::middle_names_initials(.middle_names)
                  ,
                  default_alias := 
                      std::str_lower((((.first_name)[0] ++ middle_names) ++ .last_name))
                  ,
                  conflicts := 
                      (SELECT
                          std::count(DETACHED people::Person FILTER
                              (std::str_trim(.alias, '0123456789') = default_alias)
                          )
                      )
                  ,
                  suffix := 
                      (IF (conflicts > 0) THEN <std::str>conflicts ELSE '')
              SELECT
                  (default_alias ++ suffix)
              )));
          CREATE REWRITE
              UPDATE 
              USING ((default::null_if_empty(.alias) ?? (WITH
                  middle_names := 
                      people::middle_names_initials(.middle_names)
                  ,
                  default_alias := 
                      std::str_lower((((.first_name)[0] ++ middle_names) ++ .last_name))
                  ,
                  conflicts := 
                      (SELECT
                          std::count(DETACHED people::Person FILTER
                              (std::str_trim(.alias, '0123456789') = default_alias)
                          )
                      )
                  ,
                  suffix := 
                      (IF (conflicts > 0) THEN <std::str>conflicts ELSE '')
              SELECT
                  (default_alias ++ suffix)
              )));
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(4);
      };
  };
};
