with module taxonomy,
data := <json>$data
for item in json_array_unpack(data) union (
  with anchor := <bool>item['anchor']
  insert Taxon {
      name := <str>item['name'],
      GBIF_ID := <int32>item['GBIF_ID'],
      status := <TaxonStatus>item['status'],
      parent := (
        select detached Taxon filter .GBIF_ID = <int32>item['parentID']
      ),
      rank := <Rank>item['rank'],
      authorship := <str>item['authorship'],
      anchor := anchor
    }
  unless conflict on (.name, .status)
  else (
    update Taxon set { anchor := anchor }
  )
);