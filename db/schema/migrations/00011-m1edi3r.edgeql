CREATE MIGRATION m1edi3ryzb25mbzbpxhh25fme4n45r5uo2djmtrsxn3mjwvyinguua
    ONTO m1yy6mfdzjzrbr4iosyhxz5akkoce4x2sov3lmwzabgfvp2is4xirq
{
  ALTER TYPE location::Habitat {
      ALTER LINK incompatible {
          RESET CARDINALITY;
      };
  };
};
