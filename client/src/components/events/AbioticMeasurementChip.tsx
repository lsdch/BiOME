import { AbioticMeasurement } from '@/api'
import { ComponentProps } from 'vue-component-type-helpers'
import { VChip } from 'vuetify/components'

type AbioticMeasurementChipProps = ComponentProps<typeof VChip> & {
  measurement: AbioticMeasurement
}

export function AbioticMeasurementChip({
  measurement: { value, param },
  ...props
}: AbioticMeasurementChipProps) {
  return (
    <v-chip {...props}>
      <code>
        {' '}
        {value} {param.unit}{' '}
      </code>
    </v-chip>
  )
}
