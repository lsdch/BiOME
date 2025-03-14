import { CoordinatesPrecision } from '@/api'
import { FunctionalComponent } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { VChip, VTooltip } from 'vuetify/components'

type CoordPrecisionChipProps = {
  precision: CoordinatesPrecision
} & ComponentProps<typeof VChip>

export const CoordPrecisionChip: FunctionalComponent<CoordPrecisionChipProps> = (
  { precision, ...props },
  context
) => {
  return (
    <v-tooltip>
      {{
        activator: ({ props: p }: { props: Record<string, any> }) => (
          <v-chip text={precision} {...{ ...props, ...p }} class="font-monospace" />
        ),
        default: () => CoordinatesPrecision.description[precision]
      }}
    </v-tooltip>
  )
}

export default CoordPrecisionChip
