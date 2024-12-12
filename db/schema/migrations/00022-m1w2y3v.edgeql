CREATE MIGRATION m1w2y3vrpdo43un245junhnilrkby2rmppgmxckgjrvgx7me46jueq
    ONTO m1mwk7kkzhmsguk67x3vhaiqx7q6hp75fgjkbzeslm2e6ohn7ddzrq
{
  ALTER TYPE occurrence::ExternalBioMat {
      ALTER PROPERTY is_congruent {
          USING (SELECT
              std::assert_exists((.is_homogenous AND (NOT (EXISTS (.sequences)) OR (.identification.taxon = std::assert_single(DISTINCT (.sequences.identification.taxon))))))
          );
      };
  };
  ALTER TYPE occurrence::InternalBioMat {
      ALTER PROPERTY is_congruent {
          USING (SELECT
              std::assert_exists((.is_homogenous AND (NOT (EXISTS (.specimens)) OR (.identification.taxon = std::assert_single(DISTINCT (.specimens.identification.taxon))))))
          );
      };
  };
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          *,
          required is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? ([IS occurrence::InternalBioMat].is_homogenous ?? true)),
          required is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true)),
          category := (IF (occurrence::BioMaterial IS occurrence::InternalBioMat) THEN 'Internal' ELSE (IF (occurrence::BioMaterial IS occurrence::ExternalBioMat) THEN 'External' ELSE 'Unknown'))
      }
  );
};
