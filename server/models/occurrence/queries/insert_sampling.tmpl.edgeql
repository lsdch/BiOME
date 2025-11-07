with module events,
  data := ({{ .JSON }}),
  site := ({{ .Site }}),
  select (insert events::Sampling {
    site := site,
    performed_by := (
      select people::Person
      filter .alias in <str>json_array_unpack(json_get(data, 'performed_by'))
    ),
    performed_by_groups := (
      select people::Organisation
      filter .code in <str>json_array_unpack(json_get(data,'performed_by_groups'))
    ),
    performed_on := (
      select date::from_json_with_precision(json_get(data, 'performed_on'))
    ),
    methods := (
      select SamplingMethod
      filter .code in <str>json_array_unpack(json_get(data, 'methods'))
    ),
    fixatives := (
      select samples::Fixative
      filter .code in <str>json_array_unpack(json_get(data, 'fixatives'))
    ),
    sampling_target := <SamplingTarget>(data['target']['kind']),
    target_taxa := (
      select taxonomy::Taxon
      filter .name in <str>json_array_unpack(json_get(data, 'target_taxa'))
    ),
    sampling_duration := <int32>json_get(data, 'duration'),
    comments := <str>json_get(data, 'comments'),
    habitats := (
      select sampling::Habitat
      filter .label in <str>json_array_unpack(json_get(data, 'habitats'))
    ),
    access_points := (<str>json_array_unpack(json_get(data, 'access_points')))
  })
