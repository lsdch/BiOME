CREATE MIGRATION m14wjjj7d46idowjhavxow7ipxz5pfnkklo4xms2fn2u4byipcm6jq
    ONTO m14rcaewqajbp2jjpx4pr5isb73sopxlmtjfoalndfrt7n2pi4fkza
{
  ALTER TYPE events::Sampling {
      DROP LINK occurring_taxa;
      DROP LINK reports;
      ALTER LINK samples {
          USING (.<sampling[IS occurrence::BioMaterial]);
      };
  };
};
