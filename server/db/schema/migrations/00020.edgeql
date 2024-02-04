CREATE MIGRATION m1bmp76tzguakpkyazu6ilndop5qplzmzfbnj67ap33qkt73gbcika
    ONTO m14qd5tczq32z6amxuzsrsmlu5wh7lisb7jok4grga3qjgmtin6ela
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
              USING (UPDATE
                  .meta
              SET {
                  modified := std::datetime_of_statement()
              });
      };
  };
  ALTER TYPE default::Meta {
      ALTER PROPERTY modified {
          DROP REWRITE
              UPDATE ;
          };
      };
};
