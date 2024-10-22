CREATE MIGRATION m1vz4wlfzuhbpvm64nxquge37sq7hdiiy4yb3fpaodsdqf4detf2za
    ONTO m1my2xv37bnt43aooetxzdfurgjbuuam4jkvgiw6ihydwxzmgmlqxa
{
  ALTER TYPE tokens::EmailConfirmation {
      CREATE REQUIRED LINK user_request: people::PendingUserRequest {
          ON TARGET DELETE DELETE SOURCE;
          SET REQUIRED USING (<people::PendingUserRequest>{});
          CREATE CONSTRAINT std::exclusive;
      };
      DROP PROPERTY email;
  };
};
