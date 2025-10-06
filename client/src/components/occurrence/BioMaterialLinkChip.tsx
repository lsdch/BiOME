import { BioMaterial, OccurrenceCategory } from '@/api'

/**
 * A chip displaying an occurrence category and element.
 */
export function BioMaterialLinkChip(
  { biomat: { identification, category, code } }: { biomat: BioMaterial },
  context: { attrs?: object }
) {
  return (
    <v-chip
      variant="tonal"
      text={identification.taxon.name}
      color={OccurrenceCategory.props[category].color}
      prepend-icon={OccurrenceCategory.icon(category)}
      to={{ name: 'biomat-item', params: { code: code } }}
      label
      {...context.attrs}
    />
  )
}

export default BioMaterialLinkChip
