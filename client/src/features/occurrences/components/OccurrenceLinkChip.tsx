import { BaseOccurrenceSamplingOutline } from '@/api'

/**
 * A chip displaying an occurrence category and element.
 */
export function OccurrenceLinkChip(
  { biomat: { identification, code } }: { biomat: BaseOccurrenceSamplingOutline },
  context: { attrs?: object }
) {
  return (
    <v-chip
      variant="tonal"
      text={identification.taxon.name}
      color="primary"
      prepend-icon="mdi-package-variant"
      to={{ name: 'occurrence-item', params: { code: code } }}
      label
      {...context.attrs}
    />
  )
}

export default OccurrenceLinkChip
