CREATE MIGRATION m1luh6ycpmenczqlikdltisl5kye3qlqlwjrijv2zzaqezfl5dtqca
    ONTO m1ejwxaoozwzpgwodj2lnn7xiibtrasqovv4ut72smr2qli7b7p2bq
{
  ALTER FUNCTION taxonomy::is_child(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) USING (SELECT
      ((IF (ancestor.rank = taxonomy::Rank.Kingdom) THEN (taxon.kingdom = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Phylum) THEN (taxon.phylum = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Class) THEN (taxon.class = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Order) THEN (taxon.order = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Family) THEN (taxon.family = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Genus) THEN (taxon.genus = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Species) THEN (taxon.species = ancestor) ELSE false))))))) ?? false)
  );
};
