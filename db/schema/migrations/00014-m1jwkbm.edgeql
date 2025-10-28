CREATE MIGRATION m1jwkbm73xw3duawzhgrybmkni5pxopwxzn2gsnirztre6hrdopwoa
    ONTO m1bxahepsmcbg3ewujpdijxqqucf3a4sslx2iugt5ipzqi4fkubf7a
{
  ALTER ALIAS occurrence::OccurrenceWithType USING (SELECT
      occurrence::Occurrence {
          category := (IF (.__type__.name = 'InternalBioMat') THEN occurrence::OccurrenceCategory.Internal ELSE occurrence::OccurrenceCategory.External),
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
