CREATE MIGRATION m1soaj6xcfajfl2cmlywr4rdvm6rt6swx5wmxccvdr632vhstpqvsq
    ONTO m1l7trqua4xlbezhdxal6njx6wn6qovxgxt4emr2wrk5g2ebdqp7dq
{
  ALTER TYPE admin::SecuritySettings {
      ALTER PROPERTY auth_token_lifetime {
          SET default := 600;
      };
      ALTER PROPERTY jwt_secret_key {
          CREATE CONSTRAINT std::min_len_value(32);
      };
  };
};
