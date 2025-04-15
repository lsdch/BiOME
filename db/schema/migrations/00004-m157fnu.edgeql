CREATE MIGRATION m157fnuonpk45ggpvdzbx46uwkfmkijitph2h6xuihl5cyrnwzbqba
    ONTO m1rbvgir64ynum5rwujqqgnq3mrpvtgjpteqhupd2seg7heyi7f4qa
{
  ALTER FUNCTION taxonomy::is_child(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) USING (std::assert_exists((IF (ancestor.rank = taxonomy::Rank.Kingdom) THEN (taxon.kingdom = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Phylum) THEN (taxon.phylum = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Class) THEN (taxon.class = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Order) THEN (taxon.order = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Family) THEN (taxon.family = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Genus) THEN (taxon.genus = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Species) THEN (taxon.species = ancestor) ELSE false))))))), message := ((('Failed to find taxon with name: ' ++ taxon.name) ++ ' and ancestor: ') ++ ancestor.name)));
};
