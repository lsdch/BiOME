with module people,
  data := <json>$0,
  identity := data['identity']
select (
  insert PendingUserRequest {
    email := <str>data['email'],
    identity := (
      first_name := str_trim(<str>identity['first_name']),
      last_name := str_trim(<str>identity['last_name']),
    ),
    institution := str_trim(<str>json_get(identity, 'institution')),
    motive := str_trim(<str>json_get(data, 'motive'))
  }
) { * };
