CREATE MIGRATION m1ejwxaoozwzpgwodj2lnn7xiibtrasqovv4ut72smr2qli7b7p2bq
    ONTO m157fnuonpk45ggpvdzbx46uwkfmkijitph2h6xuihl5cyrnwzbqba
{
  ALTER FUNCTION taxonomy::is_child(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) USING (std::assert_exists((IF (ancestor.rank = taxonomy::Rank.Kingdom) THEN (SELECT
      (taxon.kingdom = ancestor)
  ) ELSE (IF (ancestor.rank = taxonomy::Rank.Phylum) THEN (SELECT
      (taxon.phylum = ancestor)
  ) ELSE (IF (ancestor.rank = taxonomy::Rank.Class) THEN (SELECT
      (taxon.class = ancestor)
  ) ELSE (IF (ancestor.rank = taxonomy::Rank.Order) THEN (SELECT
      (taxon.order = ancestor)
  ) ELSE (IF (ancestor.rank = taxonomy::Rank.Family) THEN (SELECT
      (taxon.family = ancestor)
  ) ELSE (IF (ancestor.rank = taxonomy::Rank.Genus) THEN (SELECT
      (taxon.genus = ancestor)
  ) ELSE (IF (ancestor.rank = taxonomy::Rank.Species) THEN (SELECT
      (taxon.species = ancestor)
  ) ELSE (SELECT
      false
  )))))))), message := ((('Failed to find taxon with name: ' ++ taxon.name) ++ ' and ancestor: ') ++ ancestor.name)));
};
