with module taxonomy,
data := <json>$0
select (
  insert Taxon {
    name := <str>data['name'],
    code := <str>data['code'],
    GBIF_ID := <int32>data['GBIF_ID'],
    status := <TaxonStatus>data['status'],
    parent := (
      select detached Taxon filter .code = <str>data['parent']
    ),
    rank := <Rank>data['rank'],
    authorship := <str>data['authorship']
  }
) { *, parent : { * , meta: { ** }}, children : { * , meta: { ** }} };