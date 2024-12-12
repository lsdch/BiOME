CREATE MIGRATION m1x4b3ztm5kjvp5nsksrprwplv2sjbg43bn47huk54tlwe7zqffroq
    ONTO m1gczeafpaw6yy7f6x6jwlkvu7q42r4vvzd3va2w54k54bolm37tcq
{
  ALTER TYPE default::Vocabulary {
      ALTER PROPERTY label {
          DROP CONSTRAINT std::min_len_value(4);
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
};
