CREATE MIGRATION m1gzxrpzvw2qhtqu5yxfgw7t2gdh3tmo7lkenqm6ln7pfmgcrik2ea
    ONTO m14ytq2iiu6hdrp5vuhphjzouphbytocdfely5fjxozsffv6hkhyaa
{
  ALTER TYPE taxonomy::Taxon {
      CREATE CONSTRAINT std::expression ON (NOT (std::contains(.name, ' ')));
  };
};
