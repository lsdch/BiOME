CREATE MIGRATION m1rswrgjdjd3cmgkqxngwyljwscwd4fjtu5q7ke5pprqosx3l3hshq
    ONTO m1xgqekkyu3jimxajatcyzt6vslmmkzlxhq4mszwth7qzzvjnkrxaq
{
  ALTER TYPE admin::SecuritySettings {
      DROP PROPERTY account_token_lifetime;
      DROP PROPERTY auth_token_lifetime;
  };
};
