CREATE MIGRATION m1hnfypxr4yvh3gk2ozdrunnxd7ugbhnbuy7wwfsk7xshjlq6anoca
    ONTO m1wl3ufuahzlvzeebode3yyqmpqgs4vzk3dxf3mntqvduifeugocga
{
  ALTER FUNCTION taxonomy::is_in_clade(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) USING (SELECT
      (((((((((ancestor.rank = taxonomy::Rank.Kingdom) AND ((taxon.kingdom = ancestor) ?? false)) OR ((ancestor.rank = taxonomy::Rank.Phylum) AND ((taxon.phylum = ancestor) ?? false))) OR ((ancestor.rank = taxonomy::Rank.Class) AND ((taxon.class = ancestor) ?? false))) OR ((ancestor.rank = taxonomy::Rank.Order) AND ((taxon.order = ancestor) ?? false))) OR ((ancestor.rank = taxonomy::Rank.Family) AND ((taxon.family = ancestor) ?? false))) OR ((ancestor.rank = taxonomy::Rank.Genus) AND ((taxon.genus = ancestor) ?? false))) OR ((ancestor.rank = taxonomy::Rank.Species) AND ((taxon.species = ancestor) ?? false))) ?? false)
  );
};
