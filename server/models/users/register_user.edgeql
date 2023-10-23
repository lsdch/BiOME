with module people,
data := <json>$0,
select (
  insert User {
    login := <str>data['login'],
    email := <str>data['email'],
    # email_public := <str>data['email_public'],
    password := <str>data['password'],
    role := people::UserRole.Guest,
    identity := (insert Person {
      first_name := <str>data['identity']['first_name'],
      last_name := <str>data['identity']['last_name'],
      contact := <str>data['identity']['contact']
    })
  }
) {
  id, login, email, password, role, verified, identity : { * }
};