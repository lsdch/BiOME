import { H3CellWithRichness } from '@/api'

export type BindingName =
  | 'constant'
  | 'samplings'
  | 'occurrences'
  | 'speciesRichness'
  | 'genusRichness'
  | 'familyRichness'

export const bindingLabels: Record<BindingName, string> = {
  constant: 'Constant',
  samplings: 'Samplings',
  occurrences: 'Occurrences',
  speciesRichness: 'Species richness',
  genusRichness: 'Genus richness',
  familyRichness: 'Family richness'
}

export function getColorValue(d: H3CellWithRichness, binding: BindingName): number {
  switch (binding) {
    case 'constant':
      return 1
    case 'samplings':
      return d.samplings_count
    case 'occurrences':
      return d.occurrences_count
    case 'speciesRichness':
      return d.species_richness
    case 'genusRichness':
      return d.genus_richness
    case 'familyRichness':
      return d.family_richness
  }
}

export type ScaleBindingSpec = {
  binding?: BindingName
  log?: boolean
}
