CREATE MIGRATION m1pjt3ak4avjf6ktcjqsxctuk342bs2p7zz6hmnjgtaxm6gc7wslwa
    ONTO m1bmp76tzguakpkyazu6ilndop5qplzmzfbnj67ap33qkt73gbcika
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
                      modified := std::datetime_of_statement()
                  })
              );
      };
  };
};
