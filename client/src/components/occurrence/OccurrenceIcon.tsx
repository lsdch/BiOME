import { OccurrenceCategory } from '@/api'

export function OccurrenceIcon(
  { item: { category } }: { item: { category: OccurrenceCategory } },
  context: { attrs?: object }
) {
  const tooltip = `${category} occurrence`
  return (
    <v-tooltip>
      {{
        default: () => tooltip,
        activator: ({ props }: any) => (
          <v-icon
            color={category === 'Internal' ? 'primary' : 'warning'}
            icon={OccurrenceCategory.icon(category)}
            {...{ ...props, ...context.attrs }}
          />
        )
      }}
    </v-tooltip>
  )
}

export default OccurrenceIcon
