CREATE MIGRATION m1c2b4m5ojl5knshpexhq5dnc3gi7hklyszcium4ddbrvi6cg35vkq
    ONTO m1zekcs5zvbhydwi6zxeliwrmbwozwwuef3e2bsqsvtwwmvroxkfuq
{
  ALTER TYPE event::Sampling {
      ALTER LINK all_ids {
          RENAME TO occurring_taxa;
      };
  };
};
