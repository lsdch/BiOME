with module taxonomy,
taxon := (
  update Taxon
  filter .code = <str>$target
  set {
    name := <str>$name,
    code := <str>"TEST"
    # GBIF_ID := <int32>$GBIF_ID,
    # status := <TaxonStatus>$status,
    # parent := (
    #   select detached Taxon filter .GBIF_ID = <int32>$parentID
    # ),
    # rank := <Rank>$rank,
    # authorship := <str>$authorship
  }
)
select taxon { *,
  parent : { * , meta: { ** }},
  children : { * , meta: { ** }}
};