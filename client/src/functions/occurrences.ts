import { OccurrenceAtSite, TaxonInner, TaxonRank } from "@/api";

export type SampledTaxa = Record<TaxonRank, Record<string, TaxonInner>>
export function occurringTaxa(occurrences: OccurrenceAtSite[]) {
  return occurrences.reduce<SampledTaxa>((acc, occurrence) => {
    const taxon = occurrence.taxon
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
    Subspecies: {},
  })
}