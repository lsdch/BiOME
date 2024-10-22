with module people,
  data := <json>$0,
select (
  insert PendingUserRequest {
    email := <str>data['email'],
    first_name := str_trim(<str>data['first_name']),
    last_name := str_trim(<str>data['last_name']),
    institution := str_trim(<str>json_get(data, 'institution')),
    motive := str_trim(<str>json_get(data, 'motive'))
  }
) { * };
