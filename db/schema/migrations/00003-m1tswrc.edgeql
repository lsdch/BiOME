CREATE MIGRATION m1tswrc4zoyyku5yqoavbptqodgu36bc4bfjw2s4gzl34g7jd6vr2a
    ONTO m1q5yqbby6htowhaoduoxie4ob3mlwixpbyvbzc2jksjkx4gywkvxq
{
  ALTER TYPE default::Auditable {
      DROP TRIGGER update_meta;
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
};
