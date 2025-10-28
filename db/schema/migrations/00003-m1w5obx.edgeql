CREATE MIGRATION m1w5obxmcerof4vfqdom2kqv5vp6zs27ac26hzuhztlcejwz652epa
    ONTO m1m4k5casiuunmop6wytotorrbapq3eicgk66jbzmjcuged3tblpfq
{
  CREATE GLOBAL default::batch_import_id -> std::str {
      SET default := (<std::str>{});
  };
  ALTER TYPE default::Meta {
      CREATE PROPERTY batch_import_id: std::str {
          SET default := (GLOBAL default::batch_import_id);
      };
  };
};
