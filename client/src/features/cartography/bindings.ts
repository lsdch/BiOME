import { SiteWithOccurrences, TaxonRank } from '@/api'

type WeightBindingName = 'sites' | 'samplings' | 'occurrences'
type ColorBindingName = 'speciesRichness' | 'genusRichness' | 'familyRichness'

export type BindingName = WeightBindingName | ColorBindingName

export const bindingLabels: Record<BindingName, string> = {
  sites: 'Sites',
  samplings: 'Samplings',
  occurrences: 'Occurrences',
  speciesRichness: 'Species richness',
  genusRichness: 'Genus richness',
  familyRichness: 'Family richness'
}

const bindingsWeightDeck: Record<WeightBindingName, (d: SiteWithOccurrences) => number> = {
  sites: (d) => 1,
  samplings: (d) => d.samplings.length,
  occurrences: (d) => d.samplings.reduce((acc, s) => acc + s.occurrences.length, 0)
}

const bindingsColorDeck: Record<BindingName, (d: SiteWithOccurrences[]) => number> = {
  sites: (d) => d.length,
  samplings: (d) => d.reduce((acc, site) => acc + site.samplings.length, 0),
  occurrences: (d) =>
    d.reduce(
      (acc, site) => acc + site.samplings.reduce((subAcc, s) => subAcc + s.occurrences.length, 0),
      0
    ),
  speciesRichness: (d) => computeRichness(d, 'Species'),
  genusRichness: (d) => computeRichness(d, 'Genus'),
  familyRichness: (d) => computeRichness(d, 'Family')
}

function isWeightBindingName(binding: BindingName): binding is WeightBindingName {
  return ['sites', 'samplings', 'occurrences'].includes(binding)
}

function isColorBindingName(binding: BindingName): binding is ColorBindingName {
  return ['speciesRichness', 'genusRichness', 'familyRichness'].includes(binding)
}

export function getBindingFn(binding: BindingName): (d: SiteWithOccurrences[]) => number {
  return bindingsColorDeck[binding]
}

export function hexagonLayerColorBinding(spec: ScaleBindingSpec | undefined) {
  if (!spec || !spec.binding) return {}
  const { binding, log } = spec
  if (binding === 'sites' && !log) return {}
  if (isWeightBindingName(binding)) {
    const weightBinding = bindingsWeightDeck[binding]
    return !!log
      ? {
          getColorValue: (d: SiteWithOccurrences[]) =>
            Math.log(d.reduce((acc, site) => acc + weightBinding(site), 0) + 1)
        }
      : { getColorWeight: weightBinding }
  } else {
    const colorBinding = bindingsColorDeck[binding]
    console.log(binding, log)
    return !!log
      ? { getColorValue: (d: SiteWithOccurrences[]) => Math.log(colorBinding(d) + 1) }
      : { getColorValue: colorBinding }
  }
}

export type ScaleBindingSpec = {
  binding?: BindingName
  log?: boolean
}

export function computeRichness(data: SiteWithOccurrences[], rank: TaxonRank): number {
  return data.reduce((acc, { samplings }) => {
    samplings.forEach((sampling) => {
      sampling.occurrences.forEach((occurrence) => {
        const taxon = occurrence.identification.taxon
        const name = taxon[rank.toLowerCase() as keyof typeof taxon] as string | undefined
        if (name) acc.add(name)
      })
    })
    return acc
  }, new Set<string>()).size
}
