with
  sampling := ({{ .Sampling }}),
  data := {{ .JSON}},
  identification := data['identification'],
  taxon := taxonomy::taxonByName(<str>identification['taxon']),
  occurrence := (
    insert occurrence::InternalBioMat {
      code := <str>json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling.code),
      identification := (
        insert occurrence::Identification {
          taxon := taxon,
          identified_by := people::personByAlias(<str>identification['identified_by']),
          identified_on := date::from_json_with_precision(identification['identified_on']),
          confer := <bool>json_get(identification, 'confer') ?? false,
          addendum := <str>json_get(identification, 'addendum'),
        }
      ),
      sampling := sampling,
      is_type := <bool>json_get(data, 'is_type') ?? false,
    }),
  select occurrence