CREATE MIGRATION m1dz4wz3qjc563frkttebo4ahpal3a4235plmjj7sohtokgadayyxa
    ONTO m13345oyje6yqti7ta734iwitnnploqga5uzxzyeskoozibtypy7mq
{
  ALTER TYPE occurrence::BioMaterial {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE occurrence::BioMaterial {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.code THEN .code ELSE occurrence::biomat_code(.identification.taxon, .sampling)));
      };
  };
};
