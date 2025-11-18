import { Coordinates } from '@/api'
import { useClipboard } from '@vueuse/core'
import { VChip } from 'vuetify/components'

type Props = VChip['$props'] & {
  coordinates: Coordinates
}

const { copy } = useClipboard()

/**
 * A chip displaying an occurrence category and element.
 */
export function CoordinatesChip({ coordinates, ...props }: Props, context: { attrs?: object }) {
  return (
    <v-chip
      class="font-monospace"
      onClick={copy(`${coordinates.latitude}, ${coordinates.longitude}`)}
      {...props}
    >
      {coordinates.latitude}°, {coordinates.longitude}°
    </v-chip>
  )
}

export default CoordinatesChip
