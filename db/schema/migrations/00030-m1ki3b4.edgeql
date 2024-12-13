CREATE MIGRATION m1ki3b4xzbyndc644qu5cniusa7xmtuy262l6umvjwizx5dofdr2iq
    ONTO m1mg576vvgzxltwcjfiy5a4gxrax6wvvef2hbmsh2lqhzceay4gsia
{
  ALTER ALIAS seq::SequenceWithType USING (SELECT
      seq::Sequence {
          *,
          sampling := ([IS seq::AssembledSequence].sampling ?? ([IS seq::ExternalSequence].sampling ?? {})),
          identification := ([IS seq::AssembledSequence].identification ?? ([IS seq::ExternalSequence].identification ?? {})),
          category := (IF (seq::Sequence IS seq::AssembledSequence) THEN 'Internal' ELSE (IF (seq::Sequence IS seq::ExternalSequence) THEN 'External' ELSE 'Unknown'))
      }
  );
};
