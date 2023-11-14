CREATE MIGRATION m1u3xnugzn6u7xtyci4shyvaw7qzjj255pi7lihggckn2m44666cea
    ONTO m1hck7qhij72j4hkzvlvn54ydzvvrsmbnl53mxma3ndgnimxqnjfba
{
  CREATE ABSTRACT TYPE people::AccountEmailToken {
      CREATE REQUIRED PROPERTY expires: std::datetime;
      CREATE REQUIRED PROPERTY token: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  CREATE TYPE people::EmailConfirmation EXTENDING people::AccountEmailToken {
      CREATE REQUIRED LINK user: people::User {
          ON TARGET DELETE DELETE SOURCE;
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE people::PasswordReset EXTENDING people::AccountEmailToken LAST;
  ALTER TYPE people::PasswordReset {
      ALTER PROPERTY expires {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY token {
          ALTER CONSTRAINT std::exclusive {
              DROP OWNED;
          };
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
  DROP EXTENSION graphql;
};
