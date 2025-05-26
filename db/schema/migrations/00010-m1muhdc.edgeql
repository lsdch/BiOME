CREATE MIGRATION m1muhdccmtusoyegmsrid2qgzjsbx4yvcnwyiw3mwt4dw37uxb7pxa
    ONTO m1hnfypxr4yvh3gk2ozdrunnxd7ugbhnbuy7wwfsk7xshjlq6anoca
{
  ALTER TYPE taxonomy::Taxon {
      CREATE LINK lineage := (SELECT
          ((((((.kingdom UNION .phylum) UNION .class) UNION .order) UNION .family) UNION .genus) UNION .species)
      );
  };
};
