CREATE MIGRATION m1aq2yye2urlkyljkht2etmqx7ktboapunolsgh7u57r7bjpsjqqzq
    ONTO m1rswrgjdjd3cmgkqxngwyljwscwd4fjtu5q7ke5pprqosx3l3hshq
{
  ALTER TYPE admin::SecuritySettings {
      CREATE REQUIRED PROPERTY invitation_token_lifetime: std::int32 {
          SET default := 7;
          CREATE ANNOTATION std::description := 'Validity period for an account invitation in days';
          CREATE CONSTRAINT std::min_value(1);
      };
      ALTER PROPERTY refresh_token_lifetime {
          ALTER ANNOTATION std::description := 'Validity period for a session refresh token in hours';
      };
  };
  CREATE TYPE tokens::SessionRefreshToken EXTENDING tokens::Token {
      CREATE REQUIRED LINK user: people::User {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
};
