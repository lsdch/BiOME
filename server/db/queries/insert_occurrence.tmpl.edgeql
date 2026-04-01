with module occurrence,
  sampling := ({{.Sampling}}),
  data := {{.JSON}},
  identification := data['identification'],
  taxname := <str>identification['taxon'],
  taxon := assert_single((
    select taxonomy::Taxon 
    filter .name = taxname
  ) ?? (
    select taxonomy::Taxon
    filter .scientific_name = taxname
  )),
  occurrence := (
    insert Occurrence {
      sampling := sampling,
      code := <str>json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling.code),
      sources := (
        select references::DataSource
        filter .code in <str>json_array_unpack(json_get(data, 'sources'))
      ),
      # external_link := <str>json_get(data, 'external_link'),
      verbatim_identification := <str>json_get(data, 'verbatim_identification'),
      quantity := <tuple<lower:int32, upper:int32>>(
        lower := json_get(data, 'quantity', '0'),
        upper := json_get(data, 'quantity', '1')
      ),
      content_description := <str>json_get(data, 'content_description'),
      collections := assert_distinct((
        for col in json_array_unpack(json_get(data, 'collections')) union (
          <tuple<name: str, vouchers: array<str>>>(
            name := <str>col['name'],
            vouchers := <array<str>>(json_get(col, 'vouchers') ?? <json>[])
          )
        )
      )),
      comments := <str>json_get(data, 'comments'),
      published_in := (
        select distinct references::Article
        filter .code in <str>json_array_unpack(json_get(data, 'published_in'))
      ),
      identification := (
        insert occurrence::Identification {
          taxon := taxon,
          identified_by := <str>json_array_unpack(json_get(identification, 'identified_by')),
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