#
# ! This is a go template file !
#
with
  user_ID := <uuid>$0,
  token := <str>$1,
  expires := <datetime>$2
insert people::{{ . }} {
  user := (select people::User filter .id = user_ID),
  token := token,
  expires := expires
} unless conflict on (.user)
else (
  update people::{{ . }}
  set {
    token := token,
    expires := expires
  }
)