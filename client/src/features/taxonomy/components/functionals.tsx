import { TaxonRank, TaxonStatus } from '@/api'
import TaxonomyItem, { TaxonomyElement } from './TaxonomyItem.vue'

export function FTaxaNestedList(props: { items: TaxonomyElement[]; rank: TaxonRank }) {
  return (
    <div
      class="pa-0 taxa-list bg-surface"
      style={{ 'grid-column': `${props.rank == 'subgenus' ? 'species' : props.rank} / span end` }}
    >
      {props.items.map((item) => (
        <TaxonomyItem item={item} />
      ))}
    </div>
  )
}

type TaxonStatusProps = {
  icon: string
  color: string
  description: string
}

export function taxonStatusIndicatorProps(status: TaxonStatus): TaxonStatusProps {
  switch (status) {
    case 'accepted':
      return {
        icon: 'mdi-circle-medium',
        color: 'success',
        description: 'Taxon definition is accepted in GBIF database.'
      }
    case 'synonym':
      return {
        icon: 'mdi-circle-medium',
        color: 'purple',
        description:
          'Taxon definition is not accepted in GBIF database, and is considered a synonym of another taxon.'
      }
    case 'doubtful':
      return {
        icon: 'mdi-circle-medium',
        color: 'error',
        description:
          'Taxon definition is considered doubtful in GBIF database.'
      }
    case 'unreferenced':
      return {
        icon: 'mdi-circle-medium',
        color: 'primary',
        description:
          'Taxon definition is not accepted in GBIF database, but is supported by a scientific consensus.'
      }
    case 'unclassified':
      return {
        icon: 'mdi-circle-medium',
        color: 'warning',
        description:
          'Taxon definition is on-going, and is not yet supported by a scientific consensus.'
      }

    default:
      console.error(`Unhandled taxon status: ${status}`)
      return {
        icon: 'mdi-alert',
        color: 'error',
        description: 'Unhandled taxon status'
      }
  }
}

export function FTaxonStatusIndicator(props: { status: TaxonStatus }) {
  const { color, icon } = taxonStatusIndicatorProps(props.status)
  return <v-icon size="small" color={color} icon={icon} />
}
