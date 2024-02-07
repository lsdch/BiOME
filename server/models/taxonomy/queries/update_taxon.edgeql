with module taxonomy
select (
  update Taxon
  filter .id = <uuid>$id
  set {
    name := <str>$name,
    code := <str>$code,
    GBIF_ID := <int32>$GBIF_ID,
    status := <TaxonStatus>$status,
    parent := (
      select detached Taxon filter .id = <uuid>$parent
    ),
    rank := <Rank>$rank,
    authorship := <str>$authorship
  }
) { *, parent : { * , meta: { ** }}, children : { * , meta: { ** }} };