with module seq, data := <json>$0
for item in json_array_unpack(data) union (
  insert SeqDB {
    label := <str>item['label'],
    code := <str>item['code'],
    description := <str>json_get(item, 'description'),
    link_template := <str>json_get(item, 'link_template')
  }
);