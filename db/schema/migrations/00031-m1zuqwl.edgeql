CREATE MIGRATION m1zuqwlcbfunmetdcaotdj3erkxxogxfy7ksyoryxduu6q4lemitzq
    ONTO m1ki3b4xzbyndc644qu5cniusa7xmtuy262l6umvjwizx5dofdr2iq
{
  ALTER ALIAS seq::SequenceWithType USING (SELECT
      seq::Sequence {
          *,
          required sampling := std::assert_exists(([IS seq::AssembledSequence].sampling ?? [IS seq::ExternalSequence].sampling)),
          required identification := std::assert_exists(([IS seq::AssembledSequence].identification ?? [IS seq::ExternalSequence].identification)),
          required category := (IF (seq::Sequence IS seq::AssembledSequence) THEN 'Internal' ELSE (IF (seq::Sequence IS seq::ExternalSequence) THEN 'External' ELSE 'Unknown'))
      }
  );
};
