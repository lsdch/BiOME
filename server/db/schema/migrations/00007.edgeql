CREATE MIGRATION m17qorhx72yd4budctmejlk5qkhme2isbplseix7fhvgob2ku4mtqa
    ONTO m1dentl2ncfb43uknlst3ms6jsscjxo55iwjyknynb6js3ykblnvla
{
  ALTER TYPE people::Institution {
      CREATE MULTI LINK people := (.<institution[IS people::Person]);
  };
};
