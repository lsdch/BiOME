CREATE MIGRATION m1eqgei3f7r2v4eku3ka6i3bzzgiky2fqbx7ldw45nrzz3vrykrjoa
    ONTO m13qxcn5ltqbwxx6exie6bev2uxxdaxf7gk4n5jkfmb4pswi2bdfxq
{
  ALTER TYPE admin::Settings {
      CREATE REQUIRED LINK superadmin: people::User {
          SET REQUIRED USING (<people::User><std::uuid>'bdac29a2-cf40-11ef-aed1-ab96228b12ce');
      };
  };
};
