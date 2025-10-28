CREATE MIGRATION m17zlvgocpoor4xv37d6sgmdkpmls72r3otob3ccjdhwzhwyub3zca
    ONTO m1tb5acygg2z3diwwhxs4v56czuuxcygolvonkjqk6tpgbeymi32sa
{
  ALTER TYPE occurrence::Occurrence {
      DROP PROPERTY category;
  };
};
