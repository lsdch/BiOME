CREATE MIGRATION m1ssydtz6erhqk2ih6sqvh5yd7wje3ekzportluqt2ne4tzofltbea
    ONTO m1bu4c6lntjbgvenwasslcmnjzsordq233weqlo7d4mkjmze2hfurq
{
  ALTER TYPE taxonomy::Taxon {
      CREATE REQUIRED PROPERTY children_count := (SELECT
          std::count(.children)
      );
  };
};
