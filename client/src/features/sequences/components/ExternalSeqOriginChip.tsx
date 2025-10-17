import { ExtSeqOrigin } from '@/api'

export function ExternalSeqOriginChip(
  { origin }: { origin: ExtSeqOrigin },
  context: { attrs?: object }
) {
  return (
    <v-menu location="top start" origin="top start" transition="scale-transition">
      {{
        default: () => (
          <v-card
            title={origin}
            prepend-icon={ExtSeqOrigin.icon(origin)}
            subtitle="External sequence origin"
            class="bg-surface-light small-card-title"
            density="compact"
            max-width={300}
          >
            {{
              default: () => <v-card-text>{ExtSeqOrigin.description(origin)}</v-card-text>
            }}
          </v-card>
        ),
        activator: ({ props }: any) => (
          <v-chip
            text={origin}
            prepend-icon={ExtSeqOrigin.icon(origin)}
            label
            {...{ ...props, ...context.attrs }}
          />
        )
      }}
    </v-menu>
  )
}

export default ExternalSeqOriginChip
