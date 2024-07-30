with module taxonomy,
data := <json>$1
select (
  update Taxon
  filter .code = <uuid>$0
  set {
    name := <str>json_get(data,'name') ?? .name,
    code := <str>json_get(data, 'code') ?? .code,
    status := <TaxonStatus>json_get(data, 'status') ?? .status,
    parent := (
      select detached Taxon filter .code = <str>parent
    ),
    rank := <Rank>json_get(data, 'rank') ?? .rank,
    authorship := <str>json_get(data, 'authorship') ?? .authorship
  }
).id;