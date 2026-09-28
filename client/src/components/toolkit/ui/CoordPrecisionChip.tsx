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
  const unit = (precision??0) > 1000 ? 'km' : 'm'
  const displayPrecision = (precision??0) > 1000 ? (precision??0)/1000 : precision
  return (
    <v-tooltip>
      {{
        activator: ({ props: p }: { props: Record<string, any> }) => (
          <v-chip
            text={precision ? `${displayPrecision} ${unit}` : 'Unknown'}
            prepend-icon={noIcon ? undefined : 'mdi-crosshairs-question'}
            {...{ ...props, ...p }}
            class="font-monospace"
          />
        ),
        default: () => precision ? `Precision within ${displayPrecision} ${unit} of provided coordinates` : 'Unknown coordinates precision'
      }}
    </v-tooltip>
  )
}

export default CoordPrecisionChip
