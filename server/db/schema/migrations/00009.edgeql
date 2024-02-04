CREATE MIGRATION m1s5gn5r34j6ya2dfz7wphacwt6multa5bkckjmujsyaa2tmaxcvbq
    ONTO m1x7y6nlhzpybfm25kfjppgt2zzzmpxa6viwxdda3pczjhg3dqixcq
{
  ALTER TYPE default::Meta {
      ALTER PROPERTY created {
          RESET default;
          CREATE REWRITE
              INSERT 
              USING (std::datetime_of_statement());
      };
  };
};
