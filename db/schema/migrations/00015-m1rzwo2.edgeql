CREATE MIGRATION m1rzwo2z2oow3bs2e2eqgxgfzyqh5yx5mnqyohuvfarnzey4l7x6ia
    ONTO m1ao2w4urmaz4y6za2ubwo25argca4kbyxf5327iltkonuvrageyla
{
  ALTER TYPE events::Sampling {
      DROP LINK occurring_taxa;
  };
};
