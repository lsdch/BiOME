CREATE MIGRATION m1p7hwtcmxrua25olb3xg4jsg2etcvherukiud2lfbuz6ifluagbtq
    ONTO m1rllj7uzgqd2dco2gymr4e7k63za7atfp3ioygpttnxh27ijrjyea
{
  ALTER ALIAS seq::SequenceWithType USING (SELECT
      seq::Sequence {
          *,
          required sampling := std::assert_exists(([IS seq::AssembledSequence].sampling ?? [IS seq::ExternalSequence].sampling)),
          required identification := std::assert_exists(([IS seq::AssembledSequence].identification ?? [IS seq::ExternalSequence].identification))
      }
  );
};
