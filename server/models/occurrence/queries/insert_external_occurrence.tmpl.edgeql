with module occurrence,
  sampling := ({{.Sampling}}),
  data := {{.JSON}},
  identification := data['identification'],
  taxon := taxonomy::taxonByName(<str>identification['taxon']),
  occurrence := (
    insert ExternalOccurrence {
      sampling := sampling,
      code := <str>json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling.code),
      sources := (
        select references::DataSource
        filter .code in <str>json_array_unpack(json_get(data, 'sources'))
      ),
      external_link := <str>json_get(data, 'external_link'),
      original_taxon := <str>json_get(data, 'original_taxon'),
      quantity := <tuple<lower:int32, upper:int32>>(
        lower := json_get(data, 'quantity', '0'),
        upper := json_get(data, 'quantity', '1')
      ),
      content_description := <str>json_get(data, 'content_description'),
      in_collection := <str>json_get(data, 'collection'),
      item_vouchers := <str>json_array_unpack(json_get(data, 'item_vouchers')),
      comments := <str>json_get(data, 'comments'),
      published_in := (
        select distinct references::Article
        filter .code in <str>json_array_unpack(json_get(data, 'published_in'))
      ),
      identification := (
        insert occurrence::Identification {
          taxon := taxon,
          identified_by := people::personByAlias(<str>json_get(identification, 'identified_by')),
          identified_on := date::from_json_with_precision(json_get(identification, 'identified_on')),
          confer := <bool>json_get(identification, 'confer') ?? false,
          addendum := <str>json_get(identification, 'addendum'),
        }
      ),
      type_status := <occurrence::TypeStatus>json_get(data, 'type_status'),
    }
  ),
  sequences := (
    for seq_data in json_array_unpack(json_get(data, 'sequences'))
    union (
      {{ template "insert_external_sequence.tmpl.edgeql" .InsertSequenceQuery }}
    )
  ),
  select occurrence