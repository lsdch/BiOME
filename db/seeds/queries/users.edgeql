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
        alias := <str>json_get(item['identity'], 'alias') ?? {},
        comment := <str>json_get(item['identity'], 'comment') ?? {},
        institutions := (
          select Institution
          filter .code IN array_unpack(<array<str>>item['identity']['institutions'])
        )
      }
    )
  }
);