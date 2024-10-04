CREATE MIGRATION m13pw5raoumfqcpoxzla7qaj32pxzs7ytm6futncxvjkd5ilivh3uq
    ONTO m1kufxx6w67wsocldhnobu3izg4jkvaxz74vgsoydcyosfmdbewdca
{
  DROP ALIAS people::ActiveUser;
  DROP ALIAS people::InactiveUser;
  ALTER TYPE people::User {
      DROP PROPERTY is_active;
  };
};
