CREATE MIGRATION m1mwk7kkzhmsguk67x3vhaiqx7q6hp75fgjkbzeslm2e6ohn7ddzrq
    ONTO m1r4mjpijbdwjc5yjx6k4ik5cwudupidfobanqbya656uw4fdsh4ja
{
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          required is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? ([IS occurrence::InternalBioMat].is_homogenous ?? true)),
          required is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true)),
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
