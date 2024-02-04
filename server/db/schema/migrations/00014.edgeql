CREATE MIGRATION m1sydsszxcqcf7wb44niurvxxjt2yfp765ss55q3x2xxdhyx5x54dq
    ONTO m1k7ubv4nok7fv7oko2mho5vr77lxfvur327wgyv6azslsy7dpsyuq
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
              USING (SELECT
                  (INSERT
                      default::Meta
                  )
              );
      };
  };
};
