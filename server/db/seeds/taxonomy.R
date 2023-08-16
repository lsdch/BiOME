library(tidyverse)

taxonomy = read_tsv("/home/louis/Bureau/backbone-simple.txt",
  na = c("", "NA", "\\N", "{}"),
  col_names = c(
    "id", "parent", "basionym", "is_synonym", "status", "rank",
    "nom_status", "constituent_key", "origin",
    "source_taxon_key", "kingdom", "phylum", "class", "order", "family", "genus", "species", "nameId", "scientificName",
    "canonical_name", "genus_or_above", "specific_epithet", "infra_specific_epithet", "notho_type",
    "authorship", "year", "bracket_authorship", "bracket_year", "published_in", "issues"
  )
)

names(taxonomy)

which(names(taxonomy) == "published_in")

nrow(taxonomy)

head(taxonomy %>% filter(!is.na(issues)) %>% select(canonical_name, status, rank, issues))
taxonomy %>% filter(nom_status != NA)

taxonomy %>%
  filter(status == "DOUBTFUL") %>%
  nrow()
nrow(taxonomy)
taxonomy %>%
  select(origin) %>%
  unique()

taxonomy %>% filter(nameId == ori)

view = taxonomy %>% select(
  id, parent, published_in, is_synonym, source_taxon_key, status, rank, origin, genus_or_above,
  specific_epithet, infra_specific_epithet, authorship, year, bracket_authorship, bracket_year,
  kingdom, phylum, class, order, family, genus, species, canonical_name
)

view %>% filter(!is.na(published_in))

view %>%
  select(bracket_authorship) %>%
  filter(!is.na(bracket_authorship))
unique()


view %>% filter(rank == "UNRANKED")
view %>% filter(rank == "VARIETY")

view %>% filter(id == 16960141)


taxonomy %>%
  filter(canonical_name == "Proasellus") %>%
  select(notho_type)
taxonomy %>%
  filter(!is.na(notho_type)) %>%
  select(canonical_name, notho_type, rank, status)
taxonomy %>%
  filter(genus_or_above == "Proasellus") %>%
  select(id, scientificName, is_synonym, rank, authorship, bracket_authorship, canonical_name, genus_or_above, specific_epithet, infra_specific_epithet, notho_type) %>%
  filter(!is.na(notho_type))
view %>%
  filter(family == 4574, rank != "UNRANKED") %>%
  print(n = 100)

view %>%
  filter(genus == 2206245, rank != "UNRANKED") %>%
  print(n = 40)
# %>%filter(genus_or_above == "Asellus")
