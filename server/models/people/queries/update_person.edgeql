with data := <json>$1,
  institutionCodes := json_get(data, 'institutions'),
  institutions := (if exists institutionCodes then (
      select people::Institution
      filter .code in array_unpack(<array<str>>institutionCodes)
    ) else ({})
  ),
select (update people::Person filter .id = <uuid>$0
  set {
    last_name := <str>json_get(data, 'last_name') ?? .last_name,
    first_name := <str>json_get(data, 'first_name') ?? .first_name,
    contact := <str>json_get(data, 'contact') ?? .contact,
    alias := <str>json_get(data, 'alias') ?? .alias,
    comment := <str>json_get(data, "comment") ?? .comment,
    institutions := distinct institutions ?? .institutions
  }
) { *, institutions: { id, name, code }, meta: { * }};