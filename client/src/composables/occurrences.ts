import { SiteWithOccurrences, TaxonRank } from "@/api"
import { ScaleBinding } from "vue-leaflet-hexbin"

export type BindingName = 'sites' | 'occurrences' | 'samplings' | 'speciesRichness' | 'genusRichness' | 'familyRichness'

const bindings: Record<BindingName, ScaleBinding<SiteWithOccurrences>> = {
  sites: (d) => d.length,
  samplings: (d) => d.reduce((acc, { data }) => acc + data.samplings.length, 0),
  occurrences: (d) => d.reduce((acc, { data }) => acc + data.samplings.flatMap((s) => s.occurrences).length, 0),
  speciesRichness: (d) => computeRichness(d.map(({ data }) => data), 'Species'),
  genusRichness: (d) => computeRichness(d.map(({ data }) => data), 'Genus'),
  familyRichness: (d) => computeRichness(d.map(({ data }) => data), 'Family'),
}

export type ScaleBindingSpec = {
  binding?: keyof typeof bindings
  log?: boolean | 10
}

export function useScaleBinding(spec: ScaleBindingSpec | undefined): ScaleBinding<SiteWithOccurrences> | undefined {
  if (!spec) return undefined
  const { binding, log } = spec
  if (!binding) return undefined
  const scaleBinding = bindings[binding]
  switch (log) {
    case true:
      return (d) => Math.log(scaleBinding(d) + 1)
    case 10:
      return (d) => Math.log10(scaleBinding(d) + 1)
    default:
      return scaleBinding
  }
}

export function computeRichness(data: SiteWithOccurrences[], rank: TaxonRank): number {
  return data.reduce(
    (acc, { samplings }) => acc + samplings.flatMap(
      (s) => s.occurrences.filter(
        (o) => o.taxon.rank === rank || (rank === 'Species' && o.taxon.rank === 'Subspecies')
      )
    ).length, 0)
}