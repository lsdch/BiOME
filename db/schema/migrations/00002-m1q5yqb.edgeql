CREATE MIGRATION m1q5yqbby6htowhaoduoxie4ob3mlwixpbyvbzc2jksjkx4gywkvxq
    ONTO m1wo32gakw2ch5oszkvu4uuckxdr3vqvkivo6mtb7cuadhe4r6xlja
{
  ALTER TYPE default::Auditable {
      CREATE TRIGGER update_meta
          AFTER UPDATE 
          FOR EACH DO (UPDATE
              __new__.meta
          SET {
              modified := std::datetime_of_statement()
          });
      ALTER LINK meta {
          DROP REWRITE
              UPDATE ;
          };
      };
};
