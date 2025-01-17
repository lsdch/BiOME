CREATE MIGRATION m1eyhpn7bbeew4vqoz2opjo3n77ipppshey3pc5hez2degqk42t6qq
    ONTO m1mhsnbai2qljuo5j2ks4amkpli5wsjm25muhdikvivjb2jt6v5pza
{
  ALTER TYPE location::Site {
      CREATE REQUIRED PROPERTY user_defined_locality: std::bool {
          SET default := false;
      };
  };
};
