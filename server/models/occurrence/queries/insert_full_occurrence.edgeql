with
  data := {{ .JSON }},
  site := location::insert_site({{ .Site }}, true),
  samplings := (
    for sampling_data in json_array_unpack(json_get(data, 'samplings'))
    union (
      {{ .InsertSamplingQuery }}
    )
  ),
  abiotic_measurements := (
    for sampling in samplings
    for measurement_data in json_array_unpack(json_get(data, 'abiotic_measurements', sampling.code))
    union (
      insert occurrence::AbioticMeasurement {
        sampling := sampling,
        parameter := occurrence::abioticParameterByCode(<str>measurement_data['parameter']),
        value := <float32>measurement_data['value'],
        unit := occurrence::unitByCode(<str>measurement_data['unit']),
        measured_on := date::from_json_with_precision(json_get(measurement_data, 'measured_on')),
        comments := <str>json_get(measurement_data, 'comments'),
      }
    )
  ),
  internal := (
    for occ_data in json_array_unpack(json_get(data, 'samplings', 'external_occurrences'))
    union (
      {{ .InsertInternalBioMatQuery }}
    )
  ),
  for occ_data in json_array_unpack(json_get(data, 'samplings',  'internal_occurrences'))
  union (
    {{ .InsertExternalOccurrenceQuery }}
  )
