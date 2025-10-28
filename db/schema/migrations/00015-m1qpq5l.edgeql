CREATE MIGRATION m1qpq5lt4wt2e55vwvmxrtj36enpejzhywvigg5ujpftcemmxdcxxa
    ONTO m1jwkbm73xw3duawzhgrybmkni5pxopwxzn2gsnirztre6hrdopwoa
{
  ALTER TYPE occurrence::Occurrence {
      CREATE REQUIRED PROPERTY category := (std::assert_exists((IF (__source__ IS occurrence::InternalBioMat) THEN occurrence::OccurrenceCategory.Internal ELSE (IF (__source__ IS occurrence::ExternalOccurrence) THEN occurrence::OccurrenceCategory.External ELSE <occurrence::OccurrenceCategory>{})), message := (('Occurrence category for subtype ' ++ __source__.__type__.name) ++ ' is undefined')));
  };
  ALTER ALIAS occurrence::OccurrenceWithType USING (SELECT
      occurrence::Occurrence {
          required has_sequences := ([IS occurrence::InternalBioMat].has_sequences ?? ([IS occurrence::ExternalOccurrence].has_sequences ?? false)),
          internal := [IS occurrence::InternalBioMat] {
              is_homogenous,
              is_congruent,
              seq_consensus
          },
          external := [IS occurrence::ExternalOccurrence] {
              sequences,
              published_in,
              original_taxon,
              external_link,
              in_collection,
              item_vouchers,
              quantity,
              content_description
          }
      }
  );
};
