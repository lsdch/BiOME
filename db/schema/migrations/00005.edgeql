CREATE MIGRATION m1cop5taizv5aeqn6t4uzzhbicko52fssrg25ddjhrvv5kifbf5jba
    ONTO m1plhlejp6dzzfr5hwkfv3yw4lznnkezriixvaf2svvvdtdehjxf7q
{
  CREATE FUNCTION default::null_if_empty(s: std::str) -> OPTIONAL std::str USING (WITH
      trimmed := 
          std::str_trim(s)
  SELECT
      (<std::str>{} IF (std::len(trimmed) = 0) ELSE trimmed)
  );
  ALTER TYPE people::Person {
      ALTER PROPERTY contact {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.contact));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.contact));
      };
      ALTER PROPERTY middle_names {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.middle_names));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.middle_names));
      };
  };
  ALTER TYPE people::Institution {
      ALTER PROPERTY description {
          CREATE REWRITE
              INSERT 
              USING (default::null_if_empty(.description));
          CREATE REWRITE
              UPDATE 
              USING (default::null_if_empty(.description));
      };
  };
  ALTER TYPE people::Institution {
      ALTER PROPERTY name {
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
  ALTER TYPE people::Institution {
      ALTER PROPERTY name {
          DROP CONSTRAINT std::min_len_value(10);
      };
  };
};
