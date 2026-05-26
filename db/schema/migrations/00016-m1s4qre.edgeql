CREATE MIGRATION m1s4qrextmowgc36nddqgtfsijtyuukoi5nlynavkpgukuc6l72dma
    ONTO m1rzwo2z2oow3bs2e2eqgxgfzyqh5yx5mnqyohuvfarnzey4l7x6ia
{
  DROP FUNCTION people::insert_or_find_organisation(data: std::json);
  DROP FUNCTION people::insert_organisation(data: std::json);
  DROP FUNCTION people::insert_person(data: std::json);
  DROP FUNCTION people::personByAlias(alias: std::str);
  ALTER TYPE datasets::ResearchProgram {
      DROP LINK funding_agencies;
  };
  ALTER TYPE datasets::ResearchProgram {
      CREATE MULTI PROPERTY funding_agencies: std::str;
  };
  ALTER TYPE people::User {
      ALTER PROPERTY login {
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
  ALTER TYPE default::Meta {
      ALTER PROPERTY created_by {
          USING ((
              id := .created_by_user.id,
              login := .created_by_user.login,
              name := .created_by_user.identity.full_name
          ));
      };
      ALTER PROPERTY updated_by {
          USING ((
              id := .modified_by_user.id,
              login := .modified_by_user.login,
              name := .modified_by_user.identity.full_name
          ));
      };
  };
  ALTER TYPE people::Organisation {
      DROP LINK people;
      DROP PROPERTY code;
      DROP PROPERTY description;
      DROP PROPERTY kind;
      DROP PROPERTY name;
  };
  ALTER TYPE people::Person {
      DROP LINK organisations;
      DROP INDEX ON ((.alias, .last_name));
      DROP PROPERTY alias;
  };
  DROP TYPE people::Organisation;
  ALTER TYPE people::Person {
      CREATE PROPERTY organisation: std::str;
  };
  ALTER TYPE people::User {
      ALTER PROPERTY login {
          DROP CONSTRAINT std::min_len_value(5);
      };
  };
  DROP SCALAR TYPE people::OrgKind;
};
