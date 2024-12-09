CREATE MIGRATION m1muksybuie6vuxtspjeydwd6esqd2tzgws2yeh5vlrvdzn76tkiva
    ONTO m1fllsjihtlerlrpcl6f6xfoqtehis76apragc2ikfuczuss6exevq
{
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          published_in: {
              *
          },
          sampling: {
              *
          },
          identification: {
              *
          },
          meta: {
              *
          },
          category := (IF (occurrence::BioMaterial IS occurrence::InternalBioMat) THEN 'Internal' ELSE (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN 'External' ELSE 'Unknown'))
      }
  );
};
