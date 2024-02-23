with data := <json>$0,
	institutions := (
    select people::Institution
    filter .code in array_unpack(<array<str>>data['institutions'])
  )
select (insert people::Person {
  first_name := <str>data['first_name'],
  middle_names := <str>json_get(data, 'middle_names') ?? {},
  last_name := <str>data['last_name'],
  alias := <str>json_get(data, 'alias') ?? {},
  comment := <str>json_get(data, "comment") ?? {},
  institutions := distinct institutions
}) { *, institutions: { * }, meta: { * }}