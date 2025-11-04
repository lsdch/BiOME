with module occurrence,
  params := <json>$0,
  search_term := <str>json_get(params, 'search'),
  category := <str>json_get(params, 'category'),
  {{ template "taxa_variables" .TaxaFilters -}}
  with_sequences := <bool>json_get(params, 'has_sequences'),
  confer := <bool>json_get(params, 'confer'),
  status := <taxonomy::TaxonStatus>json_get(params, 'status'),
  is_type := <bool>json_get(params, 'is_type'),
  is_own := <bool>params['owned'],
items := (
  select OccurrenceWithType { * }
  filter (
    {{- if .SearchTerm }}
      (.code ilike '%%' ++ search_term ++ '%%') and
    {{ end }}
    {{- if .Category.IsSet }}
    (.category = ("occurrence::InternalBioMat" if category = "Internal" else "occurrence::ExternalOccurrence")) and
    {{ end }}
    {{ template "taxa_filters" .TaxaFilters }}
    {{- if .Status.IsSet }}
      (.identification.taxon.status = status) and
    {{ end }}
    {{- if .HasSequences.IsSet }}
      (.has_sequences = with_sequences) and
    {{ end }}
    {{- if .Confer.IsSet }}
      (.identification.confer = confer) and
    {{ end }}
    {{- if .IsType.IsSet }}
      (.is_type = is_type) and
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
    order by <str>$1
    offset <optional int64>json_get(params, 'offset')
    limit <optional int64>json_get(params, 'limit')
  ) {
    *,
    sampling: { *, site: { *, country: { * } } },
    identification: { **, identified_by: { ** } },
    meta: { * }
  },
  total_count := count(items),
};