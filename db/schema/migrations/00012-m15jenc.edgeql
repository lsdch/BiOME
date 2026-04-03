CREATE MIGRATION m15jencp3tvkehhcuqtebe6j7i2zafv7ypggtpw6fo4obclfj6yvia
    ONTO m1jyqvl5nilazvtzzvxzbswbd6md5ysrr6jekarxsocybmnrsciqzq
{
  ALTER TYPE default::CodeIdentifier {
      CREATE MULTI PROPERTY code_history: tuple<code: std::str, time: std::datetime>;
  };
};
