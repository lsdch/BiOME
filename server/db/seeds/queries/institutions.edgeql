with module people,
  data := <json>$0
for item in json_array_unpack(data) union (
  insert Institution {
    name := <str>item['name'],
    code := <str>item['code'],
    description := <str>item['description']
  }
)