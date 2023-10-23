CREATE MIGRATION m1a6icewk7tkqtbgspdwzlelkqbwotfmti4e5emljfkn2lc5ylykaa
    ONTO m1lpnacdrieqhx25v66xaa5chojohpnne7btbdka5clfu4hmvgessq
{
  ALTER TYPE people::User {
      CREATE REQUIRED PROPERTY email_public: std::bool {
          SET default := false;
      };
      ALTER PROPERTY role {
          SET REQUIRED USING (<people::UserRole>{});
      };
      ALTER PROPERTY login {
          CREATE CONSTRAINT std::max_len_value(16);
          CREATE CONSTRAINT std::min_len_value(5);
      };
  };
  ALTER TYPE people::Person {
      DROP PROPERTY second_names;
  };
};
