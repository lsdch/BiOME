CREATE MIGRATION m1sznfqu2wrzglz5x2v4mj2sfz2lpubvbxi47ev56rzfkksf6sr3fq
    ONTO m13chmki2uhe2o5nhly62o26a2htkc2b3zeupzgwtwnls4a5kqioyq
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          DROP REWRITE
              UPDATE ;
          };
      };
};
