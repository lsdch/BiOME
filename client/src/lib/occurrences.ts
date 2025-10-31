import { OccurrenceAtSite, Taxon, TaxonRank } from "@/api";

export type SampledTaxa = Record<TaxonRank, Record<string, Taxon>>
export function occurringTaxa(occurrences: OccurrenceAtSite[]) {
  return occurrences.reduce<SampledTaxa>((acc, occurrence) => {
    const taxon = occurrence.identification.taxon
    acc[taxon.rank][taxon.name] = taxon
    return acc
  }, {
    Kingdom: {},
    Phylum: {},
    Class: {},
    Order: {},
    Family: {},
    Genus: {},
    Species: {},
    Subgenus: {},
    Subspecies: {},
  })
}