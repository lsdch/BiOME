CREATE MIGRATION m14lxf2xo2pnh7vt37mjrqybcomiqx5pwarc2ic4bslrgepejzssta
    ONTO m1ze4vhskhnrwoyodypldvn7hgnwtiudu7y6wvp4lja3q7274dinsa
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
                      (.id = __subject__.id)
                  SET {
                      created := default::Meta.created
                  })
              );
      };
  };
};
