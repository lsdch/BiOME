with module people,
  data := <json>$1,
select (
  update Institution filter .code = <str>$0
  set {
    name := <str>json_get(data, 'name') ?? .name,
    code := <str>json_get(data, 'code') ?? .code,
    description := <str>json_get(data, 'description') ??.description,
    kind := <InstitutionKind>json_get(data, 'kind') ?? .kind
  }
).id;