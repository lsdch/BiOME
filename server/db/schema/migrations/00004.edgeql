CREATE MIGRATION m1emojwi2egcdl6utcg6nfsluhfzyifalybrx2uotnxa2a4o3lquaq
    ONTO m1gzxrpzvw2qhtqu5yxfgw7t2gdh3tmo7lkenqm6ln7pfmgcrik2ea
{
  ALTER TYPE taxonomy::Taxon {
      DROP CONSTRAINT std::expression ON (NOT (std::contains(.name, ' ')));
  };
};
