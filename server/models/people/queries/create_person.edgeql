with module people,
  data := <json>$0,
	institutions := (
    select Institution
    filter .code in array_unpack(<array<str>>json_get(data, 'institutions'))
  )
select (insert Person {
  first_name := <str>data['first_name'],
  last_name := <str>data['last_name'],
  contact := <str>json_get(data, 'contact') ?? {},
  alias := <str>json_get(data, 'alias') ?? {},
  comment := <str>json_get(data, "comment") ?? {},
  institutions := distinct institutions
}) { ** }