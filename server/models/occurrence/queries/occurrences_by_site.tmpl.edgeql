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
  whole_clade := <bool>json_get(params, 'whole_clade'),
  habitats := <str>json_array_unpack(json_get(params, 'habitats')),
  sampling_target_taxa_names := <str>json_array_unpack(json_get(params, 'sampling_target_taxa')),
  sampling_target_taxa := (
    if exists sampling_target_taxa_names then (
      select taxonomy::Taxon
      filter .name in sampling_target_taxa_names
    ) else <taxonomy::Taxon>{}
  ),
  sampling_target_whole_clade := <bool>json_get(params, 'sampling_target_whole_clade'),
  sampling_status := <str>json_get(params, 'include_sites'),
select location::Site {
  *,
  country: { * },
  {{ if or .Habitats .SamplingTargetTaxa }}
  samplings := (
    select .samplings
    filter (
      {{ if .Habitats }}
        all(habitats in .habitats.label) and
      {{ end }}
      {{ if .SamplingTargetTaxa }}
        {{ if .SamplingTargetWholeClade }}
          any(taxonomy::is_in_clade(.sampling_target_taxa, sampling_target_taxa))
        {{ else }}
          any(.target_taxa in sampling_target_taxa)
        {{ end }}
        and
      {{ end }}
      true
    )
  ) {{ else }}
  samplings: {{ end }} {
    id,
    date := .performed_on,
    occurring_taxa: { * },
    occurrences := (
      select .occurrences {
        id,
        code,
        category,
        identification: { confer, addendum, identified_on, taxon: { * } },
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
      ),
  }
}
filter (
# (
# 	not exists sampling_status or sampling_status = "All" or (
# 		sampling_status = "Sampled" and exists .samplings
# 	) or (
# 		sampling_status = "Occurrences" and exists .samplings.occurrences
# 	)
# ) and
  {{ if .Countries }}
    .country.code in country_codes and
  {{ end }}
  # (not exists habitats or exists .samplings) and
  # (not exists taxa or exists .samplings.occurrences) and
  # (not exists sampling_target_kinds or exists .samplings) and
  {{ if .Datasets }}
    (
      location::Site in datasets[is datasets::SiteDataset].sites
    ?? datasets[is datasets::OccurrenceDataset].sites
    ) and
  {{ end }}
  true
)
