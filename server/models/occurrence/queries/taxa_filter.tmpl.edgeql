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
  {{- if .WholeClade }}
  any({.identification.taxon
  {{- range $rank, $names := .TaxaByRank -}}
    , .identification.taxon.{{- (printf "%s" $rank) | ToLower -}}
  {{- end -}}
  } in taxa)
  {{- else }}
  .identification.taxon in taxa
  {{- end}}
  and
{{- end }}
{{end}}