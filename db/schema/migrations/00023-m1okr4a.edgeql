CREATE MIGRATION m1okr4aq5ymvtsyoxg7luhjhvhoepaw45hb5jkemoc6kh7m6kpfqwa
    ONTO m1w2y3vrpdo43un245junhnilrkby2rmppgmxckgjrvgx7me46jueq
{
  ALTER TYPE occurrence::ExternalBioMat {
      ALTER PROPERTY is_congruent {
          USING (SELECT
              std::assert_exists((.is_homogenous AND (NOT (EXISTS (.sequences)) OR (.identification.taxon IN std::assert_single(DISTINCT (.sequences.identification.taxon))))))
          );
      };
  };
  ALTER TYPE occurrence::InternalBioMat {
      ALTER PROPERTY is_congruent {
          USING (SELECT
              std::assert_exists((.is_homogenous AND (NOT (EXISTS (.specimens)) OR (.identification.taxon IN std::assert_single(DISTINCT (.specimens.identification.taxon))))))
          );
      };
  };
};
