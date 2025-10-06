import { OccurrenceCategory, OccurrenceElement } from '@/api'

function chipProps(category: OccurrenceCategory) {
  const { prependIcon, color } = OccurrenceCategory.props[category]
  return { prependIcon, color }
}

/**
 * A chip displaying an occurrence category and element.
 */
export function OccurrenceCategoryChip(
  { category, element }: { category: OccurrenceCategory; element: OccurrenceElement },
  context: { attrs?: object }
) {
  return (
    <v-chip label {...{ ...chipProps(category), ...context.attrs }}>
      {category} {OccurrenceElement.humanize(element)}
    </v-chip>
  )
}

export default OccurrenceCategoryChip
