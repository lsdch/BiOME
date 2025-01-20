CREATE MIGRATION m1x74vqa5iyqskd3eahh6ffanjidjadcaetzzg67hnhfnckx55jgqq
    ONTO m1xzlkx4c2spo237cz7ypx4rmldvfcuqaisldfwcxy5i52z3z6b66a
{
  DROP TYPE datasets::Alignment;
  ALTER TYPE datasets::MOTU {
      ALTER PROPERTY number {
          SET TYPE std::int32;
      };
  };
  CREATE TYPE datasets::SeqDataset EXTENDING datasets::AbstractDataset {
      CREATE REQUIRED MULTI LINK sequences: seq::Sequence;
  };
};
