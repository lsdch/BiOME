CREATE MIGRATION m1j7wvtnoezonp6tn2agghm5pp4gwetogas7jlohtjonzog3qouzma
    ONTO m1k5kyyo6zy73u6dtpeiirmzqqeoylvsthn4szy4ryddr4q3mvjx5q
{
  ALTER TYPE people::User {
      ALTER PROPERTY password {
          CREATE ANNOTATION std::description := 'Password hashing is done within the database, raw password must be used when creating/updating.';
      };
  };
  ALTER TYPE people::User {
      ALTER PROPERTY password {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE people::User {
      ALTER PROPERTY password {
          CREATE REWRITE
              INSERT 
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des')) ELSE .password));
      };
  };
  ALTER TYPE people::User {
      ALTER PROPERTY password {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE people::User {
      ALTER PROPERTY password {
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des')) ELSE .password));
      };
  };
};
