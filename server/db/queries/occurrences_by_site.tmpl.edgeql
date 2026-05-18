{{- /* Go Template */ -}}
with module occurrence,
  params := <json>$0,
  country_codes := <str>json_array_unpack(json_get(params, 'countries')),
  {{ template "taxa_variables" .TaxaFilters -}}
  dataset_slugs := <str>json_array_unpack(json_get(params, 'datasets')),
  datasets := (
    if exists dataset_slugs then (
      select datasets::Dataset
      filter .slug in dataset_slugs
    ) else <datasets::Dataset>{}
  ),
  site_codes := <str>json_array_unpack(json_get(params, 'site_codes')),
  whole_clade := <bool>json_get(params, 'whole_clade'),
  habitats := <str>json_array_unpack(json_get(params, 'habitats')),
  sampling_target_taxa_names := <str>json_array_unpack(json_get(params, 'sampling_target', 'taxa')),
  sampling_target_taxa := (
    if exists sampling_target_taxa_names then (
      select taxonomy::Taxon
      filter .name in sampling_target_taxa_names
    ) else <taxonomy::Taxon>{}
  ),
  sampling_target_whole_clade := <bool>json_get(params, 'sampling_target', 'whole_clade'),
  sampling_status := <str>json_get(params, 'include_sites'),
select location::Site {
  *,
  country: { * },
  samplings: {
    id,
    performed_on,
    occurrences: {
      id,
      code,
      identification: { 
        confer, 
        addendum, 
        identified_on, 
        taxon: { 
          *,
          kingdom_name := .kingdom.name,
          phylum_name := .phylum.name,
          class_name := .class.name,
          order_name := .order.name,
          family_name := .family.name ?? (.name if .rank = taxonomy::Rank.Family else <str>{}),
          genus_name := .genus.name ?? (.name if .rank = taxonomy::Rank.Genus else <str>{}),
          species_name := .species.name ?? (.name if .rank = taxonomy::Rank.Species else <str>{}),
          subspecies_name := .name if .rank = taxonomy::Rank.Subspecies else <str>{},
        } },
    }
    {{- if or .Taxa .Datasets }}
      filter (
        {{ template "taxa_filters" .TaxaFilters }}
        {{ if .Datasets }}
          any(.datasets in datasets[is datasets::OccurrenceDataset]) and
        {{ end }}
        true
      )
    {{ end}}
  } 
  {{- if or .SamplingTarget.Taxa .Habitats }}
    filter (
       {{- if .SamplingTarget.Taxa }}
        {{- if .SamplingTarget.WholeClade }}
          any({.target_taxa
            {{- range $rank, $names := .TaxaByRank -}}
              , .target_taxa.{{- (printf "%s" $rank) | ToLower -}}
            {{- end -}}
          } in sampling_target_taxa)
        {{- else }}
        .target_taxa in sampling_target_taxa
        {{- end}}
        and
      {{- end }}
      {{ if .Habitats }}
        all(habitats in .habitats.label) and
      {{ end }}
      true
    )
  {{- end }}
}
{{ if or .SiteCodes .Countries }} 
filter (
# (
# 	not exists sampling_status or sampling_status = "All" or (
# 		sampling_status = "Sampled" and exists .samplings
# 	) or (
# 		sampling_status = "Occurrences" and exists .samplings.occurrences
# 	)
# ) and
  # {{ if eq .IncludeSites  "Occurrences" }}
  #   count(.samplings.occurrences) > 0 and
  # {{ else if eq .IncludeSites  "Sampled" }}
  #   exists .samplings and
  # {{ end }}
  {{ if .SiteCodes }}
    .code in site_codes and
  {{ end }}
  {{ if .Countries }}
    .country.code in country_codes and
  {{ end }}
  # (not exists habitats or exists .samplings) and
  # (not exists taxa or exists .samplings.occurrences) and
  # (not exists sampling_target_kinds or exists .samplings) and
  # {{ if .Datasets }}
  #   (
  #     location::Site in datasets[is datasets::SiteDataset].sites
  #   ?? datasets[is datasets::OccurrenceDataset].sites
  #   ) and
  # {{ end }}
  true
)
{{ end }}
