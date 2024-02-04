CREATE MIGRATION m14qd5tczq32z6amxuzsrsmlu5wh7lisb7jok4grga3qjgmtin6ela
    ONTO m1l7pdbxh4mru2j3i6sr5wmlhgetvvhnebzo6543f6flbt5hqr35kq
{
  ALTER TYPE default::Meta {
      ALTER PROPERTY created {
          SET default := (std::datetime_of_statement());
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          SET default := (INSERT
              default::Meta
          );
          DROP REWRITE
              INSERT ;
          };
      };
};
