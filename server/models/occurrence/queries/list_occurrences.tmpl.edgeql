with module occurrence,
  params := <json>$0,
  search_term := <str>json_get(params, 'search'),
  category := <str>json_get(params, 'category'),
  taxon_name := <str>json_get(params, 'taxon'),
  taxon := (
    (select taxonomy::Taxon filter .name = taxon_name)
    if (exists taxon_name)
    else <taxonomy::Taxon>{}
  ),
  whole_clade := <bool>params['whole_clade'],
  with_sequences := <bool>json_get(params, 'has_sequences'),
  confer := <bool>json_get(params, 'confer'),
  status := <taxonomy::TaxonStatus>json_get(params, 'status'),
  is_type := <bool>json_get(params, 'is_type'),
  is_own := <bool>params['owned'],
items := (
  select OccurrenceWithType { * }
  filter (
    {{ if ne .SearchTerm "" }}
      (.code ilike '%%' ++ search_term ++ '%%') and
    {{ end }}
    {{ if .Category.IsSet }}
    (.category = ("occurrence::InternalBioMat" if category = "Internal" else "occurrence::ExternalOccurrence")) and
    {{ end }}
    {{ if .Taxon.IsSet }}
    (
      taxonomy::is_in_clade(.identification.taxon, taxon) if whole_clade
      else .identification.taxon = taxon
    ) and
    {{ end }}
    {{ if .Status.IsSet }}
      (.identification.taxon.status = status) and
    {{ end }}
    {{ if .HasSequences.IsSet }}
      (.has_sequences = with_sequences) and
    {{ end }}
    {{ if .Confer.IsSet }}
      (.identification.confer = confer) and
    {{ end }}
    {{ if .IsType.IsSet }}
      (.is_type = is_type) and
    {{ end }}
    {{ if .Filter.Owned }}
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