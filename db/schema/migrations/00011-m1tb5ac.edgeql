CREATE MIGRATION m1tb5acygg2z3diwwhxs4v56czuuxcygolvonkjqk6tpgbeymi32sa
    ONTO m1qoauopnuwhfpmsfzwrnhxrpzzwhbhh4tbpubrxqmuwsplu56vqrq
{
  ALTER TYPE default::Meta {
      CREATE INDEX ON (.batch_import_id);
  };
  ALTER TYPE occurrence::Occurrence {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
};
