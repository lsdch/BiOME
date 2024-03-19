with module location,
data := <json>$0
for item in json_array_unpack(data) union (
  insert Country {
    name := <str>item['name'],
    code := <str>item['code']
  }
  unless conflict on (.code)
)