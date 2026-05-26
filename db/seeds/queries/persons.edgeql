with module people,
  data := <json>$0
for item in json_array_unpack(data) union (
  insert Person {
    first_name := <str>item['first_name'],
    last_name := <str>item['last_name'],
    comment := <str>json_get(item, 'comment') ?? {},
    organisation := <str>json_get(item, 'organisation')
  }
);