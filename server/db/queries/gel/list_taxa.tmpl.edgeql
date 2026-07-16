with module taxonomy,
    pattern := <str>$0,
    ranks := <Rank>(array_unpack(<optional array<str>>$1)),
    status := <TaxonStatus>(<str>$2 if len(<str>$2) > 0 else <str>{}),
    parent := <optional str>$3,
    sampled_only := <bool>$4,
    synonym_of := <optional str>$5,
    n := <optional int64>$6,
select Taxon { *, meta: {*}, parent_name := .parent.name }
filter (
    {{- if .Ranks }}
        (.rank in ranks) and
    {{ end }}
    {{- if .Status }}
        (.status = status) and
    {{ end }}
    {{- if .Parent }}
        (.parent.name ilike parent) and
    {{ end }}
    {{- if .SynonymOf }}
        (synonym_of in .synonyms.name) and
    {{ end }}
    {{- if .SampledOnly }}
        exists (
            select occurrence::Occurrence
            filter .identification.taxon = Taxon
        ) and
    {{ end }}
    true
)
order by 
    {{- if .Pattern }} 
    ext::pg_trgm::word_similarity_dist(pattern, .name)
    {{- else }}
    .name
    {{- end }}
    then .rank asc then .name asc
limit n