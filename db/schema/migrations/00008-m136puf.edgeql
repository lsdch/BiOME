CREATE MIGRATION m136pufzhmard4gteawe6idtkdaxhhtnzzxnku2sbdbabt2guscrkq
    ONTO m1soaj6xcfajfl2cmlywr4rdvm6rt6swx5wmxccvdr632vhstpqvsq
{
  ALTER TYPE admin::SecuritySettings {
      ALTER PROPERTY auth_token_lifetime {
          ALTER ANNOTATION std::description := 'Validity period for an authentication token in minutes';
      };
  };
  ALTER TYPE admin::SecuritySettings {
      ALTER PROPERTY auth_token_lifetime {
          CREATE CONSTRAINT std::min_value(5);
      };
  };
  ALTER TYPE admin::SecuritySettings {
      ALTER PROPERTY auth_token_lifetime {
          SET default := 15;
          DROP CONSTRAINT std::min_value(60);
      };
  };
};
