CREATE MIGRATION m1k7ubv4nok7fv7oko2mho5vr77lxfvur327wgyv6azslsy7dpsyuq
    ONTO m13jebbrxp7su46hw6mqss7zkxpfvmkp7n7iez5j6slakjqs6sll2q
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
                      .meta
                  SET {
                      created := .created
                  })
              );
      };
  };
};
