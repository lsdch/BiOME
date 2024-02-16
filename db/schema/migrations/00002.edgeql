CREATE MIGRATION m1vkutotun7xrnliqicc7tj7px2jcomxspwj3eahgdz6hfynd6p7xq
    ONTO m1isxojxvek2iemlofoe6y4xbygej6dftk3zcsctlxhwbhf2xv4j6q
{
  ALTER TYPE people::Person {
      ALTER PROPERTY first_name {
          CREATE CONSTRAINT std::min_len_value(1);
      };
      CREATE PROPERTY middle_names: std::str {
          CREATE CONSTRAINT std::max_len_value(32);
      };
      ALTER PROPERTY full_name {
          USING (std::array_join([.first_name, .middle_names, .last_name], ' '));
      };
  };
  ALTER TYPE people::Person {
      ALTER PROPERTY first_name {
          DROP CONSTRAINT std::min_len_value(2);
      };
  };
};
