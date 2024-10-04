CREATE MIGRATION m12odyfycilfqxqe73bcrkmxmfvdvnvtcz5zldssx6qucgiluioraq
    ONTO m17jfq4b5kfsl5ttztbmjd4diur55rszloxmsuoa3us3fycms5pnxa
{
  ALTER TYPE people::AccountToken {
      DROP LINK user;
  };
  ALTER TYPE people::AccountToken RENAME TO people::Token;
  ALTER TYPE people::EmailConfirmation {
      CREATE REQUIRED PROPERTY email: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE people::PasswordReset {
      CREATE REQUIRED LINK user: people::User {
          ON TARGET DELETE DELETE SOURCE;
          SET REQUIRED USING (<people::User>{});
          CREATE DELEGATED CONSTRAINT std::exclusive;
      };
  };
};
