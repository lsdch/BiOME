CREATE MIGRATION m1hmseivjwsdw3q4tqbktdgtpokjprry2y36k3zis777zgwbsjhdla
    ONTO m1luh6ycpmenczqlikdltisl5kye3qlqlwjrijv2zzaqezfl5dtqca
{
  ALTER FUNCTION taxonomy::is_child(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) USING (SELECT
      (((((((((ancestor.rank = taxonomy::Rank.Kingdom) AND (taxon.kingdom = ancestor)) OR ((ancestor.rank = taxonomy::Rank.Phylum) AND (taxon.phylum = ancestor))) OR ((ancestor.rank = taxonomy::Rank.Class) AND (taxon.class = ancestor))) OR ((ancestor.rank = taxonomy::Rank.Order) AND (taxon.order = ancestor))) OR ((ancestor.rank = taxonomy::Rank.Family) AND (taxon.family = ancestor))) OR ((ancestor.rank = taxonomy::Rank.Genus) AND (taxon.genus = ancestor))) OR ((ancestor.rank = taxonomy::Rank.Species) AND (taxon.species = ancestor))) ?? false)
  );
};
