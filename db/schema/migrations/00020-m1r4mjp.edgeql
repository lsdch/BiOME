CREATE MIGRATION m1r4mjpijbdwjc5yjx6k4ik5cwudupidfobanqbya656uw4fdsh4ja
    ONTO m1cnpdzskteshpateawedfj34e2i5zcs22qgoqt5dqvim3kzcead4a
{
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE REQUIRED PROPERTY is_homogenous := (SELECT
          (std::count(DISTINCT (.sequences.identification.taxon)) <= 1)
      );
      CREATE REQUIRED PROPERTY is_congruent := (SELECT
          std::assert_exists((.is_homogenous AND (.identification.taxon = std::assert_single(DISTINCT (.sequences.identification.taxon)))))
      );
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE REQUIRED PROPERTY is_homogenous := (SELECT
          (std::count(DISTINCT (.specimens.identification.taxon)) <= 1)
      );
      CREATE REQUIRED PROPERTY is_congruent := (SELECT
          std::assert_exists((.is_homogenous AND (.identification.taxon = std::assert_single(DISTINCT (.specimens.identification.taxon)))))
      );
  };
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? [IS occurrence::InternalBioMat].is_homogenous),
          is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? [IS occurrence::InternalBioMat].is_congruent),
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
