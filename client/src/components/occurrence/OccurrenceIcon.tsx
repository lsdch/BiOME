import { OccurrenceCategory, OccurrenceElement } from '@/api'

export function OccurrenceIcon(
  {
    item: { category, element }
  }: { item: { category: OccurrenceCategory; element: OccurrenceElement } },
  context: { attrs?: object }
) {
  const tooltip = `${category} ${element}`
  return (
    <v-tooltip>
      {{
        default: () => tooltip,
        activator: ({ props }: any) => (
          <v-icon
            color={category === 'Internal' ? 'primary' : 'warning'}
            icon={element === 'Sequence' ? 'mdi-dna' : OccurrenceCategory.icon(category)}
            {...{ ...props, ...context.attrs }}
          />
        )
      }}
    </v-tooltip>
  )
}

export default OccurrenceIcon
