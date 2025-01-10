CREATE MIGRATION m1sd2ildlby4zjxatkwk5mxpcdwhgkgcnoom63dpt37aiuu7wzswnq
    ONTO m1jhzkpurb6uv4sdhl7qzvfvlizdcchdylrbwlawno6ij4prkhe7la
{
  DROP ALIAS occurrence::BioMaterialWithType;
  ALTER TYPE occurrence::ExternalBioMat {
      DROP LINK sequence_consensus;
  };
  ALTER TYPE occurrence::ExternalBioMat {
      ALTER PROPERTY is_homogenous {
          RENAME TO homogenous;
      };
  };
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE SINGLE LINK seq_consensus := (SELECT
          (IF .homogenous THEN std::assert_single(DISTINCT (.sequences.identification.taxon)) ELSE {})
      );
  };
  ALTER TYPE occurrence::ExternalBioMat {
      ALTER PROPERTY is_congruent {
          RENAME TO congruent;
      };
  };
  ALTER TYPE occurrence::InternalBioMat {
      DROP LINK sequence_consensus;
  };
  ALTER TYPE occurrence::InternalBioMat {
      ALTER PROPERTY is_homogenous {
          RENAME TO homogenous;
      };
  };
  ALTER TYPE occurrence::InternalBioMat {
      CREATE SINGLE LINK seq_consensus := (SELECT
          (IF .homogenous THEN std::assert_single(DISTINCT (.specimens.identification.taxon)) ELSE {})
      );
  };
  ALTER TYPE occurrence::InternalBioMat {
      ALTER PROPERTY is_congruent {
          RENAME TO congruent;
      };
  };
};
