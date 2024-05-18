CREATE MIGRATION m1xlljiu42rsst325lg6znmcq5ofu6ktsabx72r3crqq5rsgkvdpia
    ONTO m1eegiinu6cy4hjigrrtjryokro6ntnzyef2yho4mojal6hti5fsnq
{
  ALTER TYPE default::Meta {
      ALTER PROPERTY created_by {
          USING ((
              id := .created_by_user.id,
              login := .created_by_user.login,
              name := .created_by_user.identity.full_name,
              alias := .created_by_user.identity.alias
          ));
      };
      ALTER PROPERTY updated_by {
          USING ((
              id := .modified_by_user.id,
              login := .modified_by_user.login,
              name := .modified_by_user.identity.full_name,
              alias := .modified_by_user.identity.alias
          ));
      };
  };
};
