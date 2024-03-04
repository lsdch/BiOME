with module people,
  data := <json>$0,
select ( insert Institution {
  name := <str>data['name'],
  code := <str>data['code'],
  description := <str>json_get(data, 'description'),
  kind := <InstitutionKind>data['kind']
}) { *, people:{ * }, meta:{ * } };