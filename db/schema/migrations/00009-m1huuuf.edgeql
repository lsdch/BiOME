CREATE MIGRATION m1huuufuhetqpkemw73stzhyu6lujfm7abv7glzvvopwoefvi2ulnq
    ONTO m13pw5raoumfqcpoxzla7qaj32pxzs7ytm6futncxvjkd5ilivh3uq
{
  ALTER TYPE people::PendingUserRequest {
      ALTER PROPERTY email {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE people::PendingUserRequest {
      DROP PROPERTY identity;
  };
  ALTER TYPE people::PendingUserRequest {
      CREATE PROPERTY institution: std::str;
  };
};
