CREATE MIGRATION m1l7pdbxh4mru2j3i6sr5wmlhgetvvhnebzo6543f6flbt5hqr35kq
    ONTO m137q6kn3qmz5xrhz47en6obqdvsep3xzq7q2fkipre6hvjus5so3a
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          CREATE REWRITE
              INSERT 
              USING (INSERT
                  default::Meta
              );
      };
  };
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
              USING (UPDATE
                  .meta
              SET {
                  created := .created
              });
      };
  };
};
