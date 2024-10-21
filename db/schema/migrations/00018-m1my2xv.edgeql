CREATE MIGRATION m1my2xv37bnt43aooetxzdfurgjbuuam4jkvgiw6ihydwxzmgmlqxa
    ONTO m17plnintkhztel6bqic25i5wj7oqhiqnoohiycrpw3ys6w73dhj3a
{
  ALTER TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY first_name: std::str {
          SET REQUIRED USING (<std::str>.identity.first_name);
      };
  };
  ALTER TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY last_name: std::str {
          SET REQUIRED USING (<std::str>.identity.last_name);
      };
  };
  ALTER TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY full_name := (((.first_name ++ ' ') ++ .last_name));
  };
  ALTER TYPE people::PendingUserRequest {
      DROP PROPERTY identity;
  };
};
