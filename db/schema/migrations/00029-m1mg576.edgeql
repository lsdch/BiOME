CREATE MIGRATION m1mg576vvgzxltwcjfiy5a4gxrax6wvvef2hbmsh2lqhzceay4gsia
    ONTO m1f2ryaddwxtcxric6kgcdrvu6hjf57m7r4yuwedqcgwkb64yk5gya
{
  ALTER ALIAS seq::SequenceWithType USING (SELECT
      seq::Sequence {
          *,
          sampling := ([IS seq::AssembledSequence].sampling ?? ([IS seq::ExternalSequence].sampling ?? {})),
          category := (IF (seq::Sequence IS seq::AssembledSequence) THEN 'Internal' ELSE (IF (seq::Sequence IS seq::ExternalSequence) THEN 'External' ELSE 'Unknown'))
      }
  );
};
