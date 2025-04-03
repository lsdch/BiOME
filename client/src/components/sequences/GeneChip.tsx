import { Gene } from '@/api'
import { VChip } from 'vuetify/components'

export type GeneChipProps = {
  gene: Gene
} & VChip['$props']

export function GeneChip({ gene, ...chipProps }: GeneChipProps) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        activator: ({ props }: { props: any }) => (
          <v-chip text={gene.code} class="font-monospace" {...{ ...props, ...chipProps }} />
        ),
        default: () => (
          <v-card title={gene.label} class="bg-surface-light small-card-title" density="compact">
            {{
              subtitle: () => <code>{gene.code}</code>,
              prepend: () => (
                <v-badge icon="mdi-dna" color="transparent" class="mr-3">
                  <v-icon icon="mdi-tag" size="small" />
                </v-badge>
              ),
              default: () =>
                gene.description ? <v-card-text>{gene.description}</v-card-text> : null
            }}
          </v-card>
        )
      }}
    </v-menu>
  )
}

export default GeneChip
