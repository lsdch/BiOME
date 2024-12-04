CREATE MIGRATION m1ezxsuag3lmhzroybwww3qbledrfl4aeskzw7vvjnhymaff6hly6a
    ONTO m1pomweain35ndrk3guaczoj52x4kxgkddbr2ti27e2zmlfx52vzta
{
  ALTER TYPE occurrence::BioMaterial {
      ALTER PROPERTY code {
          DROP REWRITE
              INSERT ;
              DROP REWRITE
                  UPDATE ;
              };
          };
};
