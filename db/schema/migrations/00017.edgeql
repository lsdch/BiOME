CREATE MIGRATION m13chmki2uhe2o5nhly62o26a2htkc2b3zeupzgwtwnls4a5kqioyq
    ONTO m1ta6rbkwfven22mvl6szhz72a3fwl2zocpsfseg7bp2g45oh6xgha
{
  ALTER TYPE default::Meta {
      ALTER PROPERTY modified {
          DROP REWRITE
              UPDATE ;
          };
      };
};
