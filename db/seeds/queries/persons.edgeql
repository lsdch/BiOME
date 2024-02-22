with module people,
  data := <json>$0
for item in json_array_unpack(data) union (
  insert Person {
    first_name := <str>item['first_name'],
    last_name := <str>item['last_name'],
    middle_names := <str>json_get(item, 'middle_names') ?? {},
    alias := <str>json_get(item, 'alias') ?? {},
    institutions := (
      select Institution filter .code IN array_unpack(<array<str>>item['institutions'])
    )
  }
);