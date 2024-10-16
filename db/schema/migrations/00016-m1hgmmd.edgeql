CREATE MIGRATION m1hgmmdb37iq4amuzamumws4xpbji4lfeqst6o67qulxbscmlibbka
    ONTO m1aq2yye2urlkyljkht2etmqx7ktboapunolsgh7u57r7bjpsjqqzq
{
  ALTER TYPE admin::SecuritySettings {
      ALTER PROPERTY refresh_token_lifetime {
          SET default := ((24 * 30));
      };
  };
};
