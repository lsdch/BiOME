CREATE MIGRATION m14rcaewqajbp2jjpx4pr5isb73sopxlmtjfoalndfrt7n2pi4fkza
    ONTO m1iaaq45adbnihg37ukvubdwijzpes3in55xcuwnq6frnas4btx5gq
{
  CREATE ALIAS occurrence::BioMaterialWithType := (
      SELECT
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
              type := (IF (occurrence::BioMaterial IS occurrence::InternalBioMat) THEN 'Internal' ELSE (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN 'External' ELSE 'Unknown'))
          }
  );
};
