CREATE MIGRATION m12pk4hd3qqzbuejeit7y6wkgru5qrqegt4fpvravcxrfgrpfxcosa
    ONTO m1u3xnugzn6u7xtyci4shyvaw7qzjj255pi7lihggckn2m44666cea
{
  ALTER GLOBAL default::current_user_id SET default := (<std::uuid>{});
  CREATE GLOBAL default::current_user := (SELECT
      people::User {
          *,
          identity: {
              *
          }
      }
  FILTER
      (.id = GLOBAL default::current_user_id)
  );
};
