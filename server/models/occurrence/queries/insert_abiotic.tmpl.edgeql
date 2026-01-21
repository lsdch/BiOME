with
  data := {{ .JSON}},
  site := ({{ .Site }}),
  param := (select events::AbioticParameter filter .code = <str>data['param']),
insert events::AbioticMeasurement {
  site := site,
  performed_by := (
    <str>json_array_unpack(json_get(data, 'performed_by'))
  ),
  performed_on := (
    select date::from_json_with_precision(json_get(data, 'performed_on'))
  ),
  param := param,
  value := <float32>data['value']
}
