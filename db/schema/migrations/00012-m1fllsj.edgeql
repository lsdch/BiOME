CREATE MIGRATION m1fllsjihtlerlrpcl6f6xfoqtehis76apragc2ikfuczuss6exevq
    ONTO m16mzxn3hd5q3fcxsihhiamrcp3npff3bxz2ks6nbpaxy6hbnm3tua
{
  ALTER TYPE seq::AssembledSequence {
      ALTER LINK identification {
          SET OWNED;
          ALTER CONSTRAINT std::exclusive {
              SET OWNED;
          };
      };
      DROP EXTENDING occurrence::Occurrence;
      ALTER LINK identification {
          ON SOURCE DELETE DELETE TARGET;
          RESET readonly;
          RESET CARDINALITY;
          SET REQUIRED;
          SET TYPE occurrence::Identification;
      };
      ALTER LINK sampling {
          RESET readonly;
          RESET CARDINALITY;
      };
  };
};
