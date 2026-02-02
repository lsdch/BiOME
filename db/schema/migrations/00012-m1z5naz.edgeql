CREATE MIGRATION m1z5nazo2qcdec4ru2hv6356o73qxfgo27qmtcctzdepv6gnlqe6rq
    ONTO m1c7j6p22ea67pi2rsiqq5tq5ye4wtbrhjppeakpxjjqajp7lsf5cq
{
  CREATE ALIAS seq::GenericSequenceWithDetails := (
      SELECT
          seq::Sequence {
              required occurrence := std::assert_exists((SELECT
                  occurrence::Occurrence
              FILTER
                  (.id = (seq::Sequence[IS seq::AssembledSequence].specimen.biomat.id ?? seq::Sequence[IS seq::ExternalSequence].biomat.id))
              ), message := 'Failed to find upstream occurrence for seq::Sequence subtype'),
              required identification := std::assert_exists(([IS seq::AssembledSequence].specimen.biomat.identification ?? [IS seq::ExternalSequence].biomat.identification), message := 'Failed to find identification for seq::Sequence subtype'),
              required specimen_identifier := std::assert_exists(([IS seq::AssembledSequence].specimen.code ?? [IS seq::ExternalSequence].specimen_identifier)),
              internal := [IS seq::AssembledSequence] {
                  alignment_code,
                  assembled_by,
                  chromatograms,
                  specimen
              }
          }
  );
};
