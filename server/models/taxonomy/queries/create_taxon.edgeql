with module taxonomy,
data := <json>$0
select (
  insert Taxon {
    name := <str>data['name'],
    code := <str>json_get(data, 'code'),
    status := <TaxonStatus>data['status'],
    parent := (
      select detached Taxon filter .code = <str>data['parent']
    ),
    rank := <Rank>data['rank'],
    authorship := <str>json_get(data, 'authorship')
  }
) { *, meta: { * }, parent : { * , meta: { * }}, children : { * , meta: { * }} };