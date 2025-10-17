import { OccurrenceCategory } from '@/api'

function chipProps(category: OccurrenceCategory) {
  const { prependIcon, color } = OccurrenceCategory.props[category]
  return { prependIcon, color }
}

/**
 * A chip displaying an occurrence category and element.
 */
export function OccurrenceCategoryChip(
  { category }: { category: OccurrenceCategory },
  context: { attrs?: object }
) {
  return (
    <v-chip label {...{ ...chipProps(category), ...context.attrs }}>
      {category} occurrence
    </v-chip>
  )
}

export default OccurrenceCategoryChip
