CREATE MIGRATION m1n7axxe5swymlnux465vyaxe3ofjbcn4xtobsvf6uxigc67gvpm4q
    ONTO m1ajbhf6hnijsznh7bj6akrfjrrr5q3qf4zjleqr5jzqkscnsoi3ba
{
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE SINGLE LINK sequence_consensus := (SELECT
          (IF .is_homogenous THEN std::assert_single(DISTINCT (.sequences.identification.taxon)) ELSE {})
      );
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE SINGLE LINK sequence_consensus := (SELECT
          (IF .is_homogenous THEN std::assert_single(DISTINCT (.specimens.identification.taxon)) ELSE {})
      );
  };
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          required has_sequences := EXISTS (([IS occurrence::ExternalBioMat].sequences ?? [IS occurrence::InternalBioMat].specimens.sequences)),
          required is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? ([IS occurrence::InternalBioMat].is_homogenous ?? true)),
          required is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true)),
          sequence_consensus := ([IS occurrence::ExternalBioMat].sequence_consensus ?? [IS occurrence::InternalBioMat].sequence_consensus)
      }
  );
};
