{{define "taxa_variables"}}
{{- if .Taxa }}
  taxa_names := <str>json_array_unpack(json_get(params, 'taxa')),
  taxa := (
    select taxonomy::Taxon filter .name in taxa_names
  ),
  whole_clade := <bool>params['whole_clade'],
{{- end -}}
{{end}}

{{define "taxa_filters"}}
{{- if .Taxa }}
  {{ if .WholeClade }}
  (
    (.identification.taxon in taxa) or
    {{ range $rank, $names := .TaxaByRank }}
      (.identification.taxon.{{- (printf "%s" $rank) | ToLower }} in taxa) or
    {{ end }}
    false
  )
  {{ else }}
  .identification.taxon in taxa
  {{ end}}
  and
{{- end }}
{{end}}