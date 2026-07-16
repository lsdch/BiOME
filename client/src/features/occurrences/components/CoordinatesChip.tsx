import { Coordinates } from '@/api'
import { useFeedback } from '@/stores/feedback'
import { useClipboard } from '@vueuse/core'
import { VChip } from 'vuetify/components'

type Props = VChip['$props'] & {
  coordinates: Coordinates
}

const { copy } = useClipboard()

const {feedback} = useFeedback()

function copyCoordinates(coordinates: Coordinates) {
  const text = `${coordinates.latitude}, ${coordinates.longitude}`
  copy(text)
  feedback({ message: `Copied coordinates to clipboard: ${text}`, type: 'info' })
}

/**
 * A chip displaying an occurrence category and element.
 */
export function CoordinatesChip({ coordinates, ...props }: Props, context: { attrs?: object }) {
  return (
    <v-chip
      class="font-monospace"
      prepend-icon="mdi-crosshairs-gps"
      onClick={() => copyCoordinates(coordinates)}
      {...props}
    >
      {coordinates.latitude}°N, {coordinates.longitude}°E
    </v-chip>
  )
}

export default CoordinatesChip
