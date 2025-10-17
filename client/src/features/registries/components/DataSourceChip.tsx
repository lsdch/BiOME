import { DataSource } from '@/api'
import { VChip } from 'vuetify/components'

export type DataSourceChipProps = {
  source: DataSource
} & VChip['$props']

export function DataSourceChip({ source, ...chipProps }: DataSourceChipProps) {
  return (
    <v-menu
      location="top start"
      origin="top start"
      transition="scale-transition"
      open-on-focus={false}
      open-on-click
    >
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip text={source.code} {...{ ...props, ...chipProps }} />
        ),
        default: () => (
          <v-card
            title={source.label ?? 'Untitled'}
            subtitle={source.code ?? 'Unknown journal'}
            class="small-card-title bg-surface-light"
            density="compact"
            max-width={600}
          >
            <v-card-text>{source.description}</v-card-text>
            {source.url ? (
              <v-card-actions>
                <a href={source.url}>{source.url}</a>
              </v-card-actions>
            ) : null}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default DataSourceChip
