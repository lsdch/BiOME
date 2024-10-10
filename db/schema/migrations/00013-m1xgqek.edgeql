CREATE MIGRATION m1xgqekkyu3jimxajatcyzt6vslmmkzlxhq4mszwth7qzzvjnkrxaq
    ONTO m1d7my64slzvskmfk5zhhdjinqpepbhshkjplb6flbs4uxqkfmrd4q
{
  ALTER TYPE admin::SecuritySettings {
      CREATE REQUIRED PROPERTY refresh_token_lifetime: std::int32 {
          SET default := 1;
          CREATE ANNOTATION std::description := 'Validity period for a session refresh token in days';
          CREATE CONSTRAINT std::min_value(1);
      };
  };
};
