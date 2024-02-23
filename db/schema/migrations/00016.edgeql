CREATE MIGRATION m1ta6rbkwfven22mvl6szhz72a3fwl2zocpsfseg7bp2g45oh6xgha
    ONTO m1ejkaogorsq45wgejygvtptjx4zapaks7sjuwk5c775nml6fcalga
{
  ALTER TYPE taxonomy::Taxon {
      ALTER PROPERTY code {
          DROP REWRITE
              INSERT ;
              DROP REWRITE
                  UPDATE ;
              };
          };
};
