with module people, data := <json>$0
for item in json_array_unpack(data) union (
  insert User {
    login := <str>item['login'],
    email := <str>item['email'],
    password := <str>item['password'],
    role := <UserRole>item['role'],
    identity := (
      insert Person {
        first_name := <str>item['identity']['first_name'],
        last_name := <str>item['identity']['last_name'],
        comment := <str>json_get(item['identity'], 'comment') ?? {},
        organisation := <str>json_get(item['identity'], 'organisation')
      }
    )
  }
);