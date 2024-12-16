CREATE MIGRATION m12at5dpt2wnfmrg4ngtk3fwbidnzcad6iblfh5bvydafaz5cml3aa
    ONTO m1b2knfljvno2qe6xznkyc4cua7ezmdj42bdjhvvl7blzmyllqwkbq
{
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          required has_sequences := EXISTS (([IS occurrence::ExternalBioMat].sequences ?? [IS occurrence::InternalBioMat].specimens.sequences)),
          required is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? ([IS occurrence::InternalBioMat].is_homogenous ?? true)),
          required is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true))
      }
  );
  CREATE SCALAR TYPE occurrence::OccurrenceCategory EXTENDING enum<Internal, External>;
};
