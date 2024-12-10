with data := <json>$0
for item in json_array_unpack(data) union (
  insert references::Article {
    authors := <array<str>>item['authors'],
    year := <int32>item['year'],
    verbatim := <str>item['verbatim'],
    journal := <str>json_get(item, 'journal'),
    title := <str>json_get(item, 'title'),
    code := <str>json_get(item, 'code') ?? references::generate_article_code(
      <array<str>>item['authors'], <int32>item['year']
    )
  }
);