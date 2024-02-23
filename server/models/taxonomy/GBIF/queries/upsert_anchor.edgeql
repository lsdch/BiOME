with module taxonomy,
data := <json>$0,
anchor := <bool>data['anchor'],
insert Taxon {
  name := <str>data['name'],
  GBIF_ID := <int32>data['GBIF_ID'],
  status := <TaxonStatus>data['status'],
  parent := (
    select detached Taxon filter .GBIF_ID = <int32>data['parentID']
  ),
  rank := <Rank>data['rank'],
  authorship := <str>data['authorship'],
  anchor := anchor
}
unless conflict on .GBIF_ID else (
  update Taxon set {
    anchor := anchor if not .anchor else .anchor
  }
);