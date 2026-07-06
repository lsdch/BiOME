import { FunctionalComponent } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { VChip } from 'vuetify/components'

type CoordPrecisionChipProps = {
  precision?: number
  noIcon?: boolean
} & ComponentProps<typeof VChip>

export const CoordPrecisionChip: FunctionalComponent<CoordPrecisionChipProps> = (
  { precision, noIcon, ...props },
  context
) => {
  return (
    <v-tooltip>
      {{
        activator: ({ props: p }: { props: Record<string, any> }) => (
          <v-chip
            text={precision}
            prepend-icon={noIcon ? undefined : 'mdi-crosshairs-question'}
            {...{ ...props, ...p }}
            class="font-monospace"
          />
        ),
        default: () => precision ? `Within ${precision} m` : 'Unknown precision'
      }}
    </v-tooltip>
  )
}

export default CoordPrecisionChip
