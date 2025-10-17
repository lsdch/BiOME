import { GenericOccurrenceSamplingOutline, OccurrenceCategory } from '@/api'

/**
 * A chip displaying an occurrence category and element.
 */
export function OccurrenceLinkChip(
  { biomat: { identification, category, code } }: { biomat: GenericOccurrenceSamplingOutline },
  context: { attrs?: object }
) {
  return (
    <v-chip
      variant="tonal"
      text={identification.taxon.name}
      color={OccurrenceCategory.props[category].color}
      prepend-icon={OccurrenceCategory.icon(category)}
      to={{ name: 'occurrence-item', params: { code: code } }}
      label
      {...context.attrs}
    />
  )
}

export default OccurrenceLinkChip
