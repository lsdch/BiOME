CREATE MIGRATION m1bxahepsmcbg3ewujpdijxqqucf3a4sslx2iugt5ipzqi4fkubf7a
    ONTO m17zlvgocpoor4xv37d6sgmdkpmls72r3otob3ccjdhwzhwyub3zca
{
  ALTER ALIAS occurrence::OccurrenceWithType USING (SELECT
      occurrence::Occurrence {
          category := (IF (.__type__.name = 'InternalBioMat') THEN 'Internal' ELSE 'External'),
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
