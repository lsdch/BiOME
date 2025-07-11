with module people,
  data := <json>$0
for item in json_array_unpack(data) union (
  insert Person {
    first_name := <str>item['first_name'],
    last_name := <str>item['last_name'],
    alias := <str>json_get(item, 'alias') ?? {},
    comment := <str>json_get(item, 'comment') ?? {},
    organisations := (
      select Organisation filter .code IN array_unpack(<array<str>>item['organisations'])
    )
  }
);