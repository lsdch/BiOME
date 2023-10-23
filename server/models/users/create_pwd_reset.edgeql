with
  user_ID := <uuid>$0,
  token := <str>$1,
  expires := <datetime>$1
insert people::PasswordReset {
  user := (select people::User filter .id = user_ID),
  token := token,
  expires := expires
} unless conflict on (.user)
else (
  update people::PasswordReset
  filter user := (select people::User filter .id = user_ID)
  set {
    token := token,
    expires := expires
  }
)