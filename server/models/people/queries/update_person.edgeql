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
    middle_names := <str>json_get(data, 'middle_names') ?? .middle_names,
    contact := <str>json_get(data, 'contact') ?? .contact,
    alias := <str>json_get(data, 'alias') ?? .alias,
    institutions := distinct institutions ?? .institutions
  }
) { *, institutions: { id, name, code }, meta: { * }};