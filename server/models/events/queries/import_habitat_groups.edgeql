with module location,
items := json_array_unpack(<json>$0),
habitatGroups := (for item in items union (
  with habitatGroup := (insert HabitatGroup {
      label := <str>item['label'],
      exclusive_elements := <bool>json_get(item, 'exclusive_elements') ?? true
    }),
  habitatElements := (for habitat in json_array_unpack(item['elements']) union (
    insert Habitat {
      label := <str>habitat['label'],
      description := <str>json_get(habitat, 'description'),
      in_group := habitatGroup,
    })
  )
  select habitatGroup { **, elements := habitatElements}
))
for item in items union (
  with g := (update habitatGroups filter .label = <str>item['label'] set {
    depends := assert_single((
      select (Habitat union habitatGroups.elements)
      filter .label = <str>item['depends']
    ))
  })
  for habitat in json_array_unpack(item['elements']) union (
    update g.elements filter .label = <str>habitat['label'] set {
      incompatible_from := assert_single((
        select (Habitat union habitatGroups.elements)
        filter .label in <str>json_array_unpack(json_get(item, 'incompatible'))
      ))
    }
  )
)