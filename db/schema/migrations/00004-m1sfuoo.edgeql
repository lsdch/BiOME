CREATE MIGRATION m1sfuoo3wfnz2dbeb7ycflmo3pz4rwbwvgwelur4k75e2nvfdjeowa
    ONTO m1mgk54awoju224bmd7ks7gykzei3hi3uxxe5p27xvqgkynjqfpbdq
{
  ALTER TYPE occurrence::Occurrence {
      ALTER PROPERTY original_taxon {
          RENAME TO verbatim_identification;
      };
  };
};
