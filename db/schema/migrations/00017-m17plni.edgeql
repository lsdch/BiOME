CREATE MIGRATION m17plnintkhztel6bqic25i5wj7oqhiqnoohiycrpw3ys6w73dhj3a
    ONTO m1hgmmdb37iq4amuzamumws4xpbji4lfeqst6o67qulxbscmlibbka
{
  ALTER TYPE admin::EmailSettings {
      CREATE REQUIRED PROPERTY from_address: std::str {
          SET REQUIRED USING (<std::str>'ls.duchemin@univ-lyon1.fr');
      };
      CREATE REQUIRED PROPERTY from_name: std::str {
          SET REQUIRED USING (<std::str>'DarCo platform dev');
      };
  };
};
