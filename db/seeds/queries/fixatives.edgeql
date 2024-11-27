with data := <json>$0
for item in json_array_unpack(data) union (
  insert samples::Fixative {
    label := <str>item['label'],
    code := <str>item['code'],
    description := <str>json_get(item, 'description')
  }
);