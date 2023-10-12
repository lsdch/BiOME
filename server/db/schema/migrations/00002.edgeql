CREATE MIGRATION m1l5of6qlns35qygewok2pj4skrcx3ehty6f7p475r6abt3upvchpq
    ONTO m1ktyejir4iry46z6whw5qtxmm672hh6zoz2g72lzhjqgdchfxa6ya
{
  ALTER TYPE taxonomy::Taxon {
      DROP PROPERTY slug;
  };
};
