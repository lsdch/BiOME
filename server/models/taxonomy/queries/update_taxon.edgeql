with module taxonomy,
data := <json>$1
select (
  update Taxon
  filter .id = <uuid>$0
  set {
    name := <str>json_get(data,'name') ?? .name,
    code := <str>json_get(data, 'code') ?? .code,
    GBIF_ID := <int32>json_get(data, 'GBIF_ID') ?? .GBIF_ID,
    status := <TaxonStatus>json_get(data, 'status') ?? .status,
    parent := (
      select detached Taxon filter .code = <str>parent
    ),
    rank := <Rank>json_get(data, 'rank') ?? .rank,
    authorship := <str>json_get(data, 'authorship') ?? .authorship
  }
).id;