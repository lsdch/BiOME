with module occurrence,
  params := <json>$0,
  search_term := <str>json_get(params, 'search'),
  {{ template "taxa_variables" .TaxaFilters -}}
  datasets := <str>json_array_unpack(json_get(params, 'datasets')),
  year := <int32>json_get(params, 'year'),
  year_end := <int32>json_get(params, 'year_end'),
  with_sequences := <bool>json_get(params, 'has_sequences'),
  confer := <bool>json_get(params, 'confer'),
  status := <taxonomy::TaxonStatus>json_get(params, 'status'),
  ranks := <taxonomy::Rank>json_array_unpack(json_get(params, 'rank')),
  type_status := <TypeStatus>json_array_unpack(json_get(params, 'type_status')),
  is_own := <bool>params['owned'],
items := (
  select Occurrence { * }
  filter (
    {{- if .Year.IsSet }}
    (with inferred_year := datetime_get(.sampling.performed_on.date ?? .identification.identified_on.date, "year") ?? min(.published_in.year),
    select
      {{- if .YearEnd.IsNull }}
      (inferred_year >= <float64>year)
      {{- else if .YearEnd.IsSet }}
      (inferred_year >= <float64>year) and
      (inferred_year <= <float64>year_end)
      {{- else }}
      (inferred_year = <float64>year)
      {{- end}}
    ) and
    {{ end }}
    {{- if .Datasets }}
      any(.datasets.label in datasets) and
    {{ end }}
    {{- if .SearchTerm }}
    (
      (.code ilike '%%' ++ search_term ++ '%%') or
      (.sampling.site.name ilike '%%' ++ search_term ++ '%%') or
      any(.identification.identified_by ilike '%%' ++ search_term ++ '%%')
    ) and
    {{ end }}
    {{ template "taxa_filters" .TaxaFilters }}
    {{- if .Rank }}
      (.identification.taxon.rank in ranks) and
    {{ end }}
    {{- if .Status.IsSet }}
      (.identification.taxon.status = status) and
    {{ end }}
    {{- if .HasSequences.IsSet }}
      (.has_sequences = with_sequences) and
    {{ end }}
    {{- if .Confer.IsSet }}
      (.identification.confer = confer) and
    {{ end }}
    {{- if .TypeStatus.IsSet }}
      (.type_status in type_status) and
    {{ end }}
    {{- if .Filter.Owned }}
      (.meta.created_by_user = global default::current_user if (is_own and exists global default::current_user) else true)
    {{ end }}
    true
  )
),
select {
  items := (
    select items
    {{- if .Key }}
    order by {{ .OrderByString }}
    {{- end }}
    offset <optional int64>json_get(params, 'offset')
    limit <optional int64>json_get(params, 'limit')
  ) {
    *,
    sampling: { *, site: { *, country: { * } } },
    identification: { ** },
    meta: { * }
  },
  total_count := count(items),
};