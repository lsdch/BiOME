with module people,
  data := <json>$0
for item in json_array_unpack(data) union (
  insert Person {
    first_name := <str>item['first_name'],
    last_name := <str>item['last_name'],
    institution := (
      select Institution filter .code = <str>item['institution']
    )
  }
);