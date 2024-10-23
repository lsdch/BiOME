CREATE MIGRATION m1ghu4gbwltqhyi6lnd5eolqdo34akamdad2zec4vejjbelpfnineq
    ONTO m1gjaelgh4qeg5yu3g63zaisym6ba57er3lro5gzkuhilgarzx3coq
{
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
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('bf', 10)) ELSE .password));
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
              USING ((IF __specified__.password THEN ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('bf', 10)) ELSE .password));
      };
  };
};
