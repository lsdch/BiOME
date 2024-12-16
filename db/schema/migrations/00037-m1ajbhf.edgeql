CREATE MIGRATION m1ajbhf6hnijsznh7bj6akrfjrrr5q3qf4zjleqr5jzqkscnsoi3ba
    ONTO m1p7hwtcmxrua25olb3xg4jsg2etcvherukiud2lfbuz6ifluagbtq
{
  ALTER TYPE seq::Sequence {
      CREATE REQUIRED PROPERTY category := (std::assert_exists((IF (__source__ IS seq::AssembledSequence) THEN occurrence::OccurrenceCategory.Internal ELSE (IF (__source__ IS seq::ExternalSequence) THEN occurrence::OccurrenceCategory.External ELSE {})), message := (('Occurrence category for seq::Sequence subtype ' ++ __source__.__type__.name) ++ ' is undefined')));
  };
};
