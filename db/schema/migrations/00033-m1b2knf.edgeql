CREATE MIGRATION m1b2knfljvno2qe6xznkyc4cua7ezmdj42bdjhvvl7blzmyllqwkbq
    ONTO m1tuvmztmirbyofvrlqw5bv63q7dvf6gf3yj3xmnxzc22eme4jcswq
{
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          required has_sequences := EXISTS (([IS occurrence::ExternalBioMat].sequences ?? [IS occurrence::InternalBioMat].specimens.sequences)),
          required is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? ([IS occurrence::InternalBioMat].is_homogenous ?? true)),
          required is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true)),
          category := (IF (occurrence::BioMaterial IS occurrence::InternalBioMat) THEN 'Internal' ELSE (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN 'External' ELSE 'Unknown'))
      }
  );
};
