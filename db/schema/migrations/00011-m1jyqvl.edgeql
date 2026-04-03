CREATE MIGRATION m1jyqvl5nilazvtzzvxzbswbd6md5ysrr6jekarxsocybmnrsciqzq
    ONTO m1y5tt37vrqrgw37uztchtyrbyxu7re7e426y2hxt5vhdsqxsmbr6q
{
  ALTER TYPE default::CodeIdentifier {
      ALTER PROPERTY code_history {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE default::CodeIdentifier {
      DROP PROPERTY code_history;
  };
};
