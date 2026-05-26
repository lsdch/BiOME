with data := <json>$0
for item in json_array_unpack(data) union (
  insert datasets::ResearchProgram {
    label := <str>item['label'],
    code := <str>item['code'],
    description := <str>json_get(item, 'description'),
    start_year := <int32>json_get(item, 'start_year'),
    end_year := <int32>json_get(item, 'end_year'),
    managers := <str>json_array_unpack(json_get(item, 'managers')),
    funding_agencies := <str>json_array_unpack(json_get(item, 'funding_agencies'))
  }
);