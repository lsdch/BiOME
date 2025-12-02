CREATE MIGRATION m1mgk54awoju224bmd7ks7gykzei3hi3uxxe5p27xvqgkynjqfpbdq
    ONTO m1g5awc6ntbsknxzeuwwvzbqxohbyla7fkaw3sgeuvavybtike7c3a
{
  ALTER TYPE references::Collection {
      CREATE PROPERTY location: std::str;
      CREATE REQUIRED PROPERTY personal: std::bool {
          SET default := false;
      };
  };
};
