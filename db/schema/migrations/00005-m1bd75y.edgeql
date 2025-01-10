CREATE MIGRATION m1bd75yvmwip7cx7aac3h3rjn3bkjn5brf3gefmxblo26pif6y73ca
    ONTO m1giipdbqfje4e4ktps6mah36vtlan6kq3hwk6ubxyrd6uselpt6ya
{
  ALTER TYPE seq::Sequence {
      CREATE REQUIRED LINK sampling: events::Sampling {
          SET REQUIRED USING (<events::Sampling>{});
      };
  };
  ALTER ALIAS seq::SequenceWithType USING (SELECT
      seq::Sequence {
          *,
          required identification := std::assert_exists(([IS seq::AssembledSequence].identification ?? [IS seq::ExternalSequence].identification))
      }
  );
};
