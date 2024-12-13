CREATE MIGRATION m1f2ryaddwxtcxric6kgcdrvu6hjf57m7r4yuwedqcgwkb64yk5gya
    ONTO m1uqp4r4s5hj6oa5fsrjymkhjmkhh2fw6sndccn3cwyoevxsgqgexa
{
  CREATE ALIAS seq::SequenceWithType := (
      SELECT
          seq::Sequence {
              *,
              category := (IF (seq::Sequence IS seq::AssembledSequence) THEN 'Internal' ELSE (IF (seq::Sequence IS seq::ExternalSequence) THEN 'External' ELSE 'Unknown'))
          }
  );
};
