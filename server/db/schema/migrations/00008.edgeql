CREATE MIGRATION m1x7y6nlhzpybfm25kfjppgt2zzzmpxa6viwxdda3pczjhg3dqixcq
    ONTO m17qorhx72yd4budctmejlk5qkhme2isbplseix7fhvgob2ku4mtqa
{
  ALTER TYPE people::Person EXTENDING default::Auditable LAST;
  ALTER TYPE people::Institution EXTENDING default::Auditable LAST;
};
