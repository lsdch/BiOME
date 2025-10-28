with
  biomat := ({{ .Occurrence }}),
  seq_data := {{ .JSON }},
  sequence := (insert seq::ExternalSequence {
    biomat := biomat,
    code := <str>seq_data['code'],
    label := <str>seq_data['label'],
    sequence := <str>seq_data['sequence'],
    gene := seq::geneByCode(<str>seq_data['gene']),
    legacy := <tuple<id: int32, code: str, alignment_code: str>>json_get(seq_data, 'legacy'),
    specimen_identifier := <str>seq_data['specimen_identifier'],
    referenced_in := (
      for ref in json_array_unpack(json_get(seq_data, 'referenced_in'))
      union (
        insert references::SeqReference {
          db := references::dataSourceByCode(<str>ref['db']),
          accession := <str>ref['accession'],
        }
      )
    )
  }),
  select (
    if (select sequence.code) = "default_code"
    then (update sequence set {})
    else sequence
  )