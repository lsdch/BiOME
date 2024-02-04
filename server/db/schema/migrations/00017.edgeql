CREATE MIGRATION m137q6kn3qmz5xrhz47en6obqdvsep3xzq7q2fkipre6hvjus5so3a
    ONTO m14lxf2xo2pnh7vt37mjrqybcomiqx5pwarc2ic4bslrgepejzssta
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  (UPDATE
                      default::Meta
                  FILTER
                      (.id = __subject__.meta.id)
                  SET {
                      created := default::Meta.created
                  })
              );
      };
  };
};
