with module people,
  data := <json>$0,
  user := data['user'],
  identity := data['identity']
select (
  insert PendingUserRequest {
    user := (insert User {
      login := <str>user['login'],
      email := <str>user['email'],
      password := <str>user['password'],
      role := people::UserRole.Visitor,
    }),
    identity := (
      first_name := str_trim(<str>identity['first_name']),
      last_name := str_trim(<str>identity['last_name']),
      institution := str_trim(<str>identity['institution']),
    ),
    motive := str_trim(<str>data['motive'])
  }
) { *, user: { * } };
