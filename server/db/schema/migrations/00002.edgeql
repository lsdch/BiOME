CREATE MIGRATION m1apy5rrwgvslzo7myk37d63mwsidhdllsw55hxtpmebdrul6rzrsq
    ONTO m1qtgl32ctudqtqfab73qcduxdg2ckzp2ew57izwt5iejziirui2ja
{
  ALTER TYPE seq::AssembledSequence {
      CREATE CONSTRAINT std::exclusive ON ((.specimen, .is_reference)) EXCEPT (NOT (.is_reference));
  };
};
