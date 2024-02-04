CREATE MIGRATION m1ze4vhskhnrwoyodypldvn7hgnwtiudu7y6wvp4lja3q7274dinsa
    ONTO m1sydsszxcqcf7wb44niurvxxjt2yfp765ss55q3x2xxdhyx5x54dq
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
                      created := .created
                  })
              );
      };
  };
};
